package image // import "github.com/harness-community/docker-v23/integration/image"

import (
	"context"
	"strings"
	"testing"

	"github.com/harness-community/docker-v23/api/types"
	"github.com/harness-community/docker-v23/api/types/versions"
	"github.com/harness-community/docker-v23/integration/internal/container"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/skip"
)

func TestCommitInheritsEnv(t *testing.T) {
	skip.If(t, versions.LessThan(testEnv.DaemonAPIVersion(), "1.36"), "broken in earlier versions")
	skip.If(t, testEnv.DaemonInfo.OSType == "windows", "FIXME")
	defer setupTest(t)()
	client := testEnv.APIClient()
	ctx := context.Background()

	cID1 := container.Create(ctx, t, client)
	imgName := strings.ToLower(t.Name())

	commitResp1, err := client.ContainerCommit(ctx, cID1, types.ContainerCommitOptions{
		Changes:   []string{"ENV PATH=/bin"},
		Reference: imgName,
	})
	assert.NilError(t, err)

	image1, _, err := client.ImageInspectWithRaw(ctx, commitResp1.ID)
	assert.NilError(t, err)

	expectedEnv1 := []string{"PATH=/bin"}
	assert.Check(t, is.DeepEqual(expectedEnv1, image1.Config.Env))

	cID2 := container.Create(ctx, t, client, container.WithImage(image1.ID))

	commitResp2, err := client.ContainerCommit(ctx, cID2, types.ContainerCommitOptions{
		Changes:   []string{"ENV PATH=/usr/bin:$PATH"},
		Reference: imgName,
	})
	assert.NilError(t, err)

	image2, _, err := client.ImageInspectWithRaw(ctx, commitResp2.ID)
	assert.NilError(t, err)
	expectedEnv2 := []string{"PATH=/usr/bin:/bin"}
	assert.Check(t, is.DeepEqual(expectedEnv2, image2.Config.Env))
}
