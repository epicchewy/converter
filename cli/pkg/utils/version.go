package util

import "fmt"

var (
	// VERSION is set at build time from the file at api/version
	VERSION string = ""
	// GIT_REVISION is set at build time
	GIT_REVISION string = ""
)

func Version() string {
	version := VERSION

	if GIT_REVISION != "" {
		version = fmt.Sprintf("%s-%s", version, GIT_REVISION)
	}

	return version
}
