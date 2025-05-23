package images // import "github.com/harness-community/docker-v23/daemon/images"

import (
	"context"
	"errors"
	"testing"

	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/api/types/filters"
	registrytypes "github.com/harness-community/docker-v23/api/types/registry"
	"github.com/harness-community/docker-v23/errdefs"
	"github.com/harness-community/docker-v23/registry"
	"gotest.tools/v3/assert"
)

type fakeService struct {
	registry.Service
	shouldReturnError bool

	term    string
	results []registrytypes.SearchResult
}

func (s *fakeService) Search(ctx context.Context, term string, limit int, authConfig *types.AuthConfig, userAgent string, headers map[string][]string) (*registrytypes.SearchResults, error) {
	if s.shouldReturnError {
		return nil, errdefs.Unknown(errors.New("search unknown error"))
	}
	return &registrytypes.SearchResults{
		Query:      s.term,
		NumResults: len(s.results),
		Results:    s.results,
	}, nil
}

func TestSearchRegistryForImagesErrors(t *testing.T) {
	errorCases := []struct {
		filtersArgs       filters.Args
		shouldReturnError bool
		expectedError     string
	}{
		{
			expectedError:     "search unknown error",
			shouldReturnError: true,
		},
		{
			filtersArgs:   filters.NewArgs(filters.Arg("type", "custom")),
			expectedError: "invalid filter 'type'",
		},
		{
			filtersArgs:   filters.NewArgs(filters.Arg("is-automated", "invalid")),
			expectedError: "invalid filter 'is-automated=[invalid]'",
		},
		{
			filtersArgs: filters.NewArgs(
				filters.Arg("is-automated", "true"),
				filters.Arg("is-automated", "false"),
			),
			expectedError: "invalid filter 'is-automated",
		},
		{
			filtersArgs:   filters.NewArgs(filters.Arg("is-official", "invalid")),
			expectedError: "invalid filter 'is-official=[invalid]'",
		},
		{
			filtersArgs: filters.NewArgs(
				filters.Arg("is-official", "true"),
				filters.Arg("is-official", "false"),
			),
			expectedError: "invalid filter 'is-official",
		},
		{
			filtersArgs:   filters.NewArgs(filters.Arg("stars", "invalid")),
			expectedError: "invalid filter 'stars=invalid'",
		},
		{
			filtersArgs: filters.NewArgs(
				filters.Arg("stars", "1"),
				filters.Arg("stars", "invalid"),
			),
			expectedError: "invalid filter 'stars=invalid'",
		},
	}
	for _, tc := range errorCases {
		tc := tc
		t.Run(tc.expectedError, func(t *testing.T) {
			daemon := &ImageService{
				registryService: &fakeService{
					shouldReturnError: tc.shouldReturnError,
				},
			}
			_, err := daemon.SearchRegistryForImages(context.Background(), tc.filtersArgs, "term", 0, nil, map[string][]string{})
			assert.ErrorContains(t, err, tc.expectedError)
			if tc.shouldReturnError {
				assert.Check(t, errdefs.IsUnknown(err), "got: %T: %v", err, err)
				return
			}
			assert.Check(t, errdefs.IsInvalidParameter(err), "got: %T: %v", err, err)
		})
	}
}

func TestSearchRegistryForImages(t *testing.T) {
	term := "term"
	successCases := []struct {
		name            string
		filtersArgs     filters.Args
		registryResults []registrytypes.SearchResult
		expectedResults []registrytypes.SearchResult
	}{
		{
			name:            "empty results",
			registryResults: []registrytypes.SearchResult{},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			name: "no filter",
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
				},
			},
		},
		{
			name:        "is-automated=true, no results",
			filtersArgs: filters.NewArgs(filters.Arg("is-automated", "true")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
				},
			},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			name:        "is-automated=true",
			filtersArgs: filters.NewArgs(filters.Arg("is-automated", "true")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsAutomated: true,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsAutomated: true,
				},
			},
		},
		{
			name:        "is-automated=false, no results",
			filtersArgs: filters.NewArgs(filters.Arg("is-automated", "false")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsAutomated: true,
				},
			},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			name:        "is-automated=false",
			filtersArgs: filters.NewArgs(filters.Arg("is-automated", "false")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsAutomated: false,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsAutomated: false,
				},
			},
		},
		{
			name:        "is-official=true, no results",
			filtersArgs: filters.NewArgs(filters.Arg("is-official", "true")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
				},
			},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			name:        "is-official=true",
			filtersArgs: filters.NewArgs(filters.Arg("is-official", "true")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsOfficial:  true,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsOfficial:  true,
				},
			},
		},
		{
			name:        "is-official=false, no results",
			filtersArgs: filters.NewArgs(filters.Arg("is-official", "false")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsOfficial:  true,
				},
			},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			name:        "is-official=false",
			filtersArgs: filters.NewArgs(filters.Arg("is-official", "false")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsOfficial:  false,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsOfficial:  false,
				},
			},
		},
		{
			name:        "stars=0",
			filtersArgs: filters.NewArgs(filters.Arg("stars", "0")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					StarCount:   0,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					StarCount:   0,
				},
			},
		},
		{
			name:        "stars=0, no results",
			filtersArgs: filters.NewArgs(filters.Arg("stars", "1")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					StarCount:   0,
				},
			},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			name:        "stars=1",
			filtersArgs: filters.NewArgs(filters.Arg("stars", "1")),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name0",
					Description: "description0",
					StarCount:   0,
				},
				{
					Name:        "name1",
					Description: "description1",
					StarCount:   1,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name1",
					Description: "description1",
					StarCount:   1,
				},
			},
		},
		{
			name: "stars=1, is-official=true, is-automated=true",
			filtersArgs: filters.NewArgs(
				filters.Arg("stars", "1"),
				filters.Arg("is-official", "true"),
				filters.Arg("is-automated", "true"),
			),
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name0",
					Description: "description0",
					StarCount:   0,
					IsOfficial:  true,
					IsAutomated: true,
				},
				{
					Name:        "name1",
					Description: "description1",
					StarCount:   1,
					IsOfficial:  false,
					IsAutomated: true,
				},
				{
					Name:        "name2",
					Description: "description2",
					StarCount:   1,
					IsOfficial:  true,
					IsAutomated: false,
				},
				{
					Name:        "name3",
					Description: "description3",
					StarCount:   2,
					IsOfficial:  true,
					IsAutomated: true,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name3",
					Description: "description3",
					StarCount:   2,
					IsOfficial:  true,
					IsAutomated: true,
				},
			},
		},
	}
	for _, tc := range successCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			daemon := &ImageService{
				registryService: &fakeService{
					term:    term,
					results: tc.registryResults,
				},
			}
			results, err := daemon.SearchRegistryForImages(context.Background(), tc.filtersArgs, term, 0, nil, map[string][]string{})
			assert.NilError(t, err)
			assert.Equal(t, results.Query, term)
			assert.Equal(t, results.NumResults, len(tc.expectedResults))
			assert.DeepEqual(t, results.Results, tc.expectedResults)
		})
	}
}
