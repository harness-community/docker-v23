package dockerfile // import "github.com/harness-community/docker-v23/builder/dockerfile"

func defaultShellForOS(os string) []string {
	if os == "linux" {
		return []string{"/bin/sh", "-c"}
	}
	return []string{"cmd", "/S", "/C"}
}
