package images // import "github.com/harness-community/docker-v23/daemon/images"

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/containerd/containerd/leases"
	"github.com/containerd/containerd/namespaces"
	dist "github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/distribution"
	progressutils "github.com/harness-community/docker-v23/distribution/utils"
	"github.com/harness-community/docker-v23/errdefs"
	"github.com/harness-community/docker-v23/pkg/progress"
	"github.com/harness-community/docker-v23/pkg/streamformatter"
	"github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// PullImage initiates a pull operation. image is the repository name to pull, and
// tag may be either empty, or indicate a specific tag to pull.
func (i *ImageService) PullImage(ctx context.Context, image, tag string, platform *specs.Platform, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	start := time.Now()
	// Special case: "pull -a" may send an image name with a
	// trailing :. This is ugly, but let's not break API
	// compatibility.
	image = strings.TrimSuffix(image, ":")

	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return errdefs.InvalidParameter(err)
	}

	if tag != "" {
		// The "tag" could actually be a digest.
		var dgst digest.Digest
		dgst, err = digest.Parse(tag)
		if err == nil {
			ref, err = reference.WithDigest(reference.TrimNamed(ref), dgst)
		} else {
			ref, err = reference.WithTag(ref, tag)
		}
		if err != nil {
			return errdefs.InvalidParameter(err)
		}
	}

	err = i.pullImageWithReference(ctx, ref, platform, metaHeaders, authConfig, outStream)
	imageActions.WithValues("pull").UpdateSince(start)
	if err != nil {
		return err
	}

	if platform != nil {
		// If --platform was specified, check that the image we pulled matches
		// the expected platform. This check is for situations where the image
		// is a single-arch image, in which case (for backward compatibility),
		// we allow the image to have a non-matching architecture. The code
		// below checks for this situation, and returns a warning to the client,
		// as well as logging it to the daemon logs.
		img, err := i.GetImage(image, platform)

		// Note that this is a special case where GetImage returns both an image
		// and an error: https://github.com/harness-community/docker-v23/blob/v20.10.7/daemon/images/image.go#L175-L183
		if errdefs.IsNotFound(err) && img != nil {
			po := streamformatter.NewJSONProgressOutput(outStream, false)
			progress.Messagef(po, "", `WARNING: %s`, err.Error())
			logrus.WithError(err).WithField("image", image).Warn("ignoring platform mismatch on single-arch image")
		}
	}

	return nil
}

func (i *ImageService) pullImageWithReference(ctx context.Context, ref reference.Named, platform *specs.Platform, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	// Include a buffer so that slow client connections don't affect
	// transfer performance.
	progressChan := make(chan progress.Progress, 100)

	writesDone := make(chan struct{})

	ctx, cancelFunc := context.WithCancel(ctx)

	go func() {
		progressutils.WriteDistributionProgress(cancelFunc, outStream, progressChan)
		close(writesDone)
	}()

	ctx = namespaces.WithNamespace(ctx, i.contentNamespace)
	// Take out a temporary lease for everything that gets persisted to the content store.
	// Before the lease is cancelled, any content we want to keep should have it's own lease applied.
	ctx, done, err := tempLease(ctx, i.leases)
	if err != nil {
		return err
	}
	defer done(ctx)

	cs := &contentStoreForPull{
		ContentStore: i.content,
		leases:       i.leases,
	}
	imageStore := &imageStoreForPull{
		ImageConfigStore: distribution.NewImageConfigStoreFromStore(i.imageStore),
		ingested:         cs,
		leases:           i.leases,
	}

	imagePullConfig := &distribution.ImagePullConfig{
		Config: distribution.Config{
			MetaHeaders:      metaHeaders,
			AuthConfig:       authConfig,
			ProgressOutput:   progress.ChanOutput(progressChan),
			RegistryService:  i.registryService,
			ImageEventLogger: i.LogImageEvent,
			MetadataStore:    i.distributionMetadataStore,
			ImageStore:       imageStore,
			ReferenceStore:   i.referenceStore,
		},
		DownloadManager: i.downloadManager,
		Platform:        platform,
	}

	err = distribution.Pull(ctx, ref, imagePullConfig, cs)
	close(progressChan)
	<-writesDone
	return err
}

// GetRepository returns a repository from the registry.
func (i *ImageService) GetRepository(ctx context.Context, ref reference.Named, authConfig *types.AuthConfig) (dist.Repository, error) {
	return distribution.GetRepository(ctx, ref, &distribution.ImagePullConfig{
		Config: distribution.Config{
			AuthConfig:      authConfig,
			RegistryService: i.registryService,
		},
	})
}

func tempLease(ctx context.Context, mgr leases.Manager) (context.Context, func(context.Context) error, error) {
	nop := func(context.Context) error { return nil }
	_, ok := leases.FromContext(ctx)
	if ok {
		return ctx, nop, nil
	}

	// Use an expiration that ensures the lease is cleaned up at some point if there is a crash, SIGKILL, etc.
	opts := []leases.Opt{
		leases.WithRandomID(),
		leases.WithExpiration(24 * time.Hour),
		leases.WithLabels(map[string]string{
			"moby.lease/temporary": time.Now().UTC().Format(time.RFC3339Nano),
		}),
	}
	l, err := mgr.Create(ctx, opts...)
	if err != nil {
		return ctx, nop, errors.Wrap(err, "error creating temporary lease")
	}

	ctx = leases.WithLease(ctx, l.ID)
	return ctx, func(ctx context.Context) error {
		return mgr.Delete(ctx, l)
	}, nil
}
