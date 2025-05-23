package daemon // import "github.com/harness-community/docker-v23/daemon"

import (
	"context"
	"strconv"
	"time"

	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/container"
	libcontainerdtypes "github.com/harness-community/docker-v23/libcontainerd/types"
	"github.com/harness-community/docker-v23/restartmanager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (daemon *Daemon) setStateCounter(c *container.Container) {
	switch c.StateString() {
	case "paused":
		stateCtr.set(c.ID, "paused")
	case "running":
		stateCtr.set(c.ID, "running")
	default:
		stateCtr.set(c.ID, "stopped")
	}
}

func (daemon *Daemon) handleContainerExit(c *container.Container, e *libcontainerdtypes.EventInfo) error {
	c.Lock()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ec, et, err := daemon.containerd.DeleteTask(ctx, c.ID)
	cancel()
	if err != nil {
		logrus.WithError(err).WithField("container", c.ID).Warnf("failed to delete container from containerd")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	c.StreamConfig.Wait(ctx)
	cancel()

	c.Reset(false)

	exitStatus := container.ExitStatus{
		ExitCode: int(ec),
		ExitedAt: et,
	}
	if e != nil {
		exitStatus.ExitCode = int(e.ExitCode)
		exitStatus.ExitedAt = e.ExitedAt
		exitStatus.OOMKilled = e.OOMKilled
		if e.Error != nil {
			c.SetError(e.Error)
		}
	}

	daemonShutdown := daemon.IsShuttingDown()
	execDuration := time.Since(c.StartedAt)
	restart, wait, err := c.RestartManager().ShouldRestart(ec, daemonShutdown || c.HasBeenManuallyStopped, execDuration)
	if err != nil {
		logrus.WithError(err).
			WithField("container", c.ID).
			WithField("restartCount", c.RestartCount).
			WithField("exitStatus", exitStatus).
			WithField("daemonShuttingDown", daemonShutdown).
			WithField("hasBeenManuallyStopped", c.HasBeenManuallyStopped).
			WithField("execDuration", execDuration).
			Warn("ShouldRestart failed, container will not be restarted")
		restart = false
	}

	// cancel healthcheck here, they will be automatically
	// restarted if/when the container is started again
	daemon.stopHealthchecks(c)
	attributes := map[string]string{
		"exitCode": strconv.Itoa(int(ec)),
	}
	daemon.Cleanup(c)

	if restart {
		c.RestartCount++
		logrus.WithField("container", c.ID).
			WithField("restartCount", c.RestartCount).
			WithField("exitStatus", exitStatus).
			WithField("manualRestart", c.HasBeenManuallyRestarted).
			Debug("Restarting container")
		c.SetRestarting(&exitStatus)
	} else {
		c.SetStopped(&exitStatus)
		if !c.HasBeenManuallyRestarted {
			defer daemon.autoRemove(c)
		}
	}
	defer c.Unlock() // needs to be called before autoRemove

	daemon.setStateCounter(c)
	cpErr := c.CheckpointTo(daemon.containersReplica)

	daemon.LogContainerEventWithAttributes(c, "die", attributes)

	if restart {
		go func() {
			err := <-wait
			if err == nil {
				// daemon.netController is initialized when daemon is restoring containers.
				// But containerStart will use daemon.netController segment.
				// So to avoid panic at startup process, here must wait util daemon restore done.
				daemon.waitForStartupDone()
				if err = daemon.containerStart(c, "", "", false); err != nil {
					logrus.Debugf("failed to restart container: %+v", err)
				}
			}
			if err != nil {
				c.Lock()
				c.SetStopped(&exitStatus)
				daemon.setStateCounter(c)
				c.CheckpointTo(daemon.containersReplica)
				c.Unlock()
				defer daemon.autoRemove(c)
				if err != restartmanager.ErrRestartCanceled {
					logrus.Errorf("restartmanger wait error: %+v", err)
				}
			}
		}()
	}

	return cpErr
}

// ProcessEvent is called by libcontainerd whenever an event occurs
func (daemon *Daemon) ProcessEvent(id string, e libcontainerdtypes.EventType, ei libcontainerdtypes.EventInfo) error {
	c, err := daemon.GetContainer(id)
	if err != nil {
		return errors.Wrapf(err, "could not find container %s", id)
	}

	switch e {
	case libcontainerdtypes.EventOOM:
		// StateOOM is Linux specific and should never be hit on Windows
		if isWindows {
			return errors.New("received StateOOM from libcontainerd on Windows. This should never happen")
		}

		c.Lock()
		defer c.Unlock()
		daemon.updateHealthMonitor(c)
		if err := c.CheckpointTo(daemon.containersReplica); err != nil {
			return err
		}

		daemon.LogContainerEvent(c, "oom")
	case libcontainerdtypes.EventExit:
		if ei.ProcessID == libcontainerdtypes.InitProcessName {
			return daemon.handleContainerExit(c, &ei)
		}

		exitCode := 127
		if execConfig := c.ExecCommands.Get(ei.ProcessID); execConfig != nil {
			ec := int(ei.ExitCode)
			execConfig.Lock()
			defer execConfig.Unlock()
			execConfig.ExitCode = &ec
			execConfig.Running = false

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			execConfig.StreamConfig.Wait(ctx)
			cancel()

			if err := execConfig.CloseStreams(); err != nil {
				logrus.Errorf("failed to cleanup exec %s streams: %s", c.ID, err)
			}

			// remove the exec command from the container's store only and not the
			// daemon's store so that the exec command can be inspected.
			c.ExecCommands.Delete(execConfig.ID, execConfig.Pid)

			exitCode = ec
		}
		attributes := map[string]string{
			"execID":   ei.ProcessID,
			"exitCode": strconv.Itoa(exitCode),
		}
		daemon.LogContainerEventWithAttributes(c, "exec_die", attributes)
	case libcontainerdtypes.EventStart:
		c.Lock()
		defer c.Unlock()

		// This is here to handle start not generated by docker
		if !c.Running {
			c.SetRunning(int(ei.Pid), false)
			c.HasBeenManuallyStopped = false
			c.HasBeenStartedBefore = true
			daemon.setStateCounter(c)

			daemon.initHealthMonitor(c)

			if err := c.CheckpointTo(daemon.containersReplica); err != nil {
				return err
			}
			daemon.LogContainerEvent(c, "start")
		}

	case libcontainerdtypes.EventPaused:
		c.Lock()
		defer c.Unlock()

		if !c.Paused {
			c.Paused = true
			daemon.setStateCounter(c)
			daemon.updateHealthMonitor(c)
			if err := c.CheckpointTo(daemon.containersReplica); err != nil {
				return err
			}
			daemon.LogContainerEvent(c, "pause")
		}
	case libcontainerdtypes.EventResumed:
		c.Lock()
		defer c.Unlock()

		if c.Paused {
			c.Paused = false
			daemon.setStateCounter(c)
			daemon.updateHealthMonitor(c)

			if err := c.CheckpointTo(daemon.containersReplica); err != nil {
				return err
			}
			daemon.LogContainerEvent(c, "unpause")
		}
	}
	return nil
}

func (daemon *Daemon) autoRemove(c *container.Container) {
	c.Lock()
	ar := c.HostConfig.AutoRemove
	c.Unlock()
	if !ar {
		return
	}

	err := daemon.ContainerRm(c.ID, &types.ContainerRmConfig{ForceRemove: true, RemoveVolume: true})
	if err == nil {
		return
	}
	if c := daemon.containers.Get(c.ID); c == nil {
		return
	}

	logrus.WithError(err).WithField("container", c.ID).Error("error removing container")
}
