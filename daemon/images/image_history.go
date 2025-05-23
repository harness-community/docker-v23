package images // import "github.com/harness-community/docker-v23/daemon/images"

import (
	"fmt"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/harness-community/docker-v23/api/types/image"
	"github.com/harness-community/docker-v23/layer"
)

// ImageHistory returns a slice of ImageHistory structures for the specified image
// name by walking the image lineage.
func (i *ImageService) ImageHistory(name string) ([]*image.HistoryResponseItem, error) {
	start := time.Now()
	img, err := i.GetImage(name, nil)
	if err != nil {
		return nil, err
	}

	history := []*image.HistoryResponseItem{}

	layerCounter := 0
	rootFS := *img.RootFS
	rootFS.DiffIDs = nil

	for _, h := range img.History {
		var layerSize int64

		if !h.EmptyLayer {
			if len(img.RootFS.DiffIDs) <= layerCounter {
				return nil, fmt.Errorf("too many non-empty layers in History section")
			}
			rootFS.Append(img.RootFS.DiffIDs[layerCounter])
			l, err := i.layerStore.Get(rootFS.ChainID())
			if err != nil {
				return nil, err
			}
			layerSize = l.DiffSize()
			layer.ReleaseAndLog(i.layerStore, l)
			layerCounter++
		}

		history = append([]*image.HistoryResponseItem{{
			ID:        "<missing>",
			Created:   h.Created.Unix(),
			CreatedBy: h.CreatedBy,
			Comment:   h.Comment,
			Size:      layerSize,
		}}, history...)
	}

	// Fill in image IDs and tags
	histImg := img
	id := img.ID()
	for _, h := range history {
		h.ID = id.String()

		var tags []string
		for _, r := range i.referenceStore.References(id.Digest()) {
			if _, ok := r.(reference.NamedTagged); ok {
				tags = append(tags, reference.FamiliarString(r))
			}
		}

		h.Tags = tags

		id = histImg.Parent
		if id == "" {
			break
		}
		histImg, err = i.GetImage(id.String(), nil)
		if err != nil {
			break
		}
	}
	imageActions.WithValues("history").UpdateSince(start)
	return history, nil
}
