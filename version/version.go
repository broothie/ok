package version

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var version string

func String() string {
	return strings.TrimSpace(version)
}
