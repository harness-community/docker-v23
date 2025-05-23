package zfs // import "github.com/harness-community/docker-v23/daemon/graphdriver/zfs"

import (
	"github.com/harness-community/docker-v23/daemon/graphdriver"
	"github.com/sirupsen/logrus"
)

func checkRootdirFs(rootDir string) error {
	fsMagic, err := graphdriver.GetFSMagic(rootDir)
	if err != nil {
		return err
	}
	backingFS := "unknown"
	if fsName, ok := graphdriver.FsNames[fsMagic]; ok {
		backingFS = fsName
	}

	if fsMagic != graphdriver.FsMagicZfs {
		logrus.WithField("root", rootDir).WithField("backingFS", backingFS).WithField("storage-driver", "zfs").Error("No zfs dataset found for root")
		return graphdriver.ErrPrerequisites
	}

	return nil
}

func getMountpoint(id string) string {
	return id
}
