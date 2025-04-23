package urlutil // import "github.com/DevanshMathur19/docker-v23/builder/remotecontext/urlutil"

import "testing"

var (
	gitUrls = []string{
		"git://github.com/DevanshMathur19/docker-v23",
		"git@github.com:docker/docker.git",
		"git@bitbucket.org:atlassianlabs/atlassian-docker.git",
		"https://github.com/DevanshMathur19/docker-v23.git",
		"http://github.com/DevanshMathur19/docker-v23.git",
		"http://github.com/DevanshMathur19/docker-v23.git#branch",
		"http://github.com/DevanshMathur19/docker-v23.git#:dir",
	}
	incompleteGitUrls = []string{
		"github.com/DevanshMathur19/docker-v23",
	}
	invalidGitUrls = []string{
		"http://github.com/DevanshMathur19/docker-v23.git:#branch",
	}
)

func TestIsGIT(t *testing.T) {
	for _, url := range gitUrls {
		if !IsGitURL(url) {
			t.Fatalf("%q should be detected as valid Git url", url)
		}
	}

	for _, url := range incompleteGitUrls {
		if !IsGitURL(url) {
			t.Fatalf("%q should be detected as valid Git url", url)
		}
	}

	for _, url := range invalidGitUrls {
		if IsGitURL(url) {
			t.Fatalf("%q should not be detected as valid Git prefix", url)
		}
	}
}
