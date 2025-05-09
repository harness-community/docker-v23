package containerimage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/distribution/reference"
	distref "github.com/docker/distribution/reference"
	"github.com/harness-community/docker-v23/image"
	"github.com/harness-community/docker-v23/layer"
	"github.com/moby/buildkit/exporter"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/util/compression"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

const (
	keyImageName      = "name"
	keyBuildInfo      = "buildinfo"
	keyBuildInfoAttrs = "buildinfo-attrs"
)

// Differ can make a moby layer from a snapshot
type Differ interface {
	EnsureLayer(ctx context.Context, key string) ([]layer.DiffID, error)
}

type ImageTagger interface {
	TagImageWithReference(imageID image.ID, newTag reference.Named) error
}

// Opt defines a struct for creating new exporter
type Opt struct {
	ImageStore  image.Store
	Differ      Differ
	ImageTagger ImageTagger
}

type imageExporter struct {
	opt Opt
}

// New creates a new moby imagestore exporter
func New(opt Opt) (exporter.Exporter, error) {
	im := &imageExporter{opt: opt}
	return im, nil
}

func (e *imageExporter) Resolve(ctx context.Context, opt map[string]string) (exporter.ExporterInstance, error) {
	i := &imageExporterInstance{
		imageExporter: e,
		buildInfo:     true,
	}
	for k, v := range opt {
		switch k {
		case keyImageName:
			for _, v := range strings.Split(v, ",") {
				ref, err := distref.ParseNormalizedNamed(v)
				if err != nil {
					return nil, err
				}
				i.targetNames = append(i.targetNames, ref)
			}
		case keyBuildInfo:
			if v == "" {
				i.buildInfo = true
				continue
			}
			b, err := strconv.ParseBool(v)
			if err != nil {
				return nil, errors.Wrapf(err, "non-bool value specified for %s", k)
			}
			i.buildInfo = b
		case keyBuildInfoAttrs:
			if v == "" {
				i.buildInfoAttrs = false
				continue
			}
			b, err := strconv.ParseBool(v)
			if err != nil {
				return nil, errors.Wrapf(err, "non-bool value specified for %s", k)
			}
			i.buildInfoAttrs = b
		default:
			if i.meta == nil {
				i.meta = make(map[string][]byte)
			}
			i.meta[k] = []byte(v)
		}
	}
	return i, nil
}

type imageExporterInstance struct {
	*imageExporter
	targetNames    []distref.Named
	meta           map[string][]byte
	buildInfo      bool
	buildInfoAttrs bool
}

func (e *imageExporterInstance) Name() string {
	return "exporting to image"
}

func (e *imageExporterInstance) Config() exporter.Config {
	return exporter.Config{
		Compression: compression.Config{
			Type: compression.Default,
		},
	}
}

func (e *imageExporterInstance) Export(ctx context.Context, inp exporter.Source, sessionID string) (map[string]string, error) {
	if len(inp.Refs) > 1 {
		return nil, fmt.Errorf("exporting multiple references to image store is currently unsupported")
	}

	ref := inp.Ref
	if ref != nil && len(inp.Refs) == 1 {
		return nil, fmt.Errorf("invalid exporter input: Ref and Refs are mutually exclusive")
	}

	// only one loop
	for _, v := range inp.Refs {
		ref = v
	}

	var config []byte
	var buildInfo []byte
	switch len(inp.Refs) {
	case 0:
		config = inp.Metadata[exptypes.ExporterImageConfigKey]
		if v, ok := inp.Metadata[exptypes.ExporterBuildInfo]; ok {
			buildInfo = v
		}
	case 1:
		platformsBytes, ok := inp.Metadata[exptypes.ExporterPlatformsKey]
		if !ok {
			return nil, fmt.Errorf("cannot export image, missing platforms mapping")
		}
		var p exptypes.Platforms
		if err := json.Unmarshal(platformsBytes, &p); err != nil {
			return nil, errors.Wrapf(err, "failed to parse platforms passed to exporter")
		}
		if len(p.Platforms) != len(inp.Refs) {
			return nil, errors.Errorf("number of platforms does not match references %d %d", len(p.Platforms), len(inp.Refs))
		}
		config = inp.Metadata[fmt.Sprintf("%s/%s", exptypes.ExporterImageConfigKey, p.Platforms[0].ID)]
		if v, ok := inp.Metadata[fmt.Sprintf("%s/%s", exptypes.ExporterBuildInfo, p.Platforms[0].ID)]; ok {
			buildInfo = v
		}
	}

	var diffs []digest.Digest
	if ref != nil {
		layersDone := oneOffProgress(ctx, "exporting layers")

		if err := ref.Finalize(ctx); err != nil {
			return nil, layersDone(err)
		}

		if err := ref.Extract(ctx, nil); err != nil {
			return nil, err
		}

		diffIDs, err := e.opt.Differ.EnsureLayer(ctx, ref.ID())
		if err != nil {
			return nil, layersDone(err)
		}

		diffs = make([]digest.Digest, len(diffIDs))
		for i := range diffIDs {
			diffs[i] = digest.Digest(diffIDs[i])
		}

		_ = layersDone(nil)
	}

	if len(config) == 0 {
		var err error
		config, err = emptyImageConfig()
		if err != nil {
			return nil, err
		}
	}

	history, err := parseHistoryFromConfig(config)
	if err != nil {
		return nil, err
	}

	diffs, history = normalizeLayersAndHistory(diffs, history, ref)

	config, err = patchImageConfig(config, diffs, history, inp.Metadata[exptypes.ExporterInlineCache], buildInfo)
	if err != nil {
		return nil, err
	}

	configDigest := digest.FromBytes(config)

	configDone := oneOffProgress(ctx, fmt.Sprintf("writing image %s", configDigest))
	id, err := e.opt.ImageStore.Create(config)
	if err != nil {
		return nil, configDone(err)
	}
	_ = configDone(nil)

	if e.opt.ImageTagger != nil {
		for _, targetName := range e.targetNames {
			tagDone := oneOffProgress(ctx, "naming to "+targetName.String())
			if err := e.opt.ImageTagger.TagImageWithReference(image.ID(digest.Digest(id)), targetName); err != nil {
				return nil, tagDone(err)
			}
			_ = tagDone(nil)
		}
	}

	return map[string]string{
		exptypes.ExporterImageConfigDigestKey: configDigest.String(),
		exptypes.ExporterImageDigestKey:       id.String(),
	}, nil
}
