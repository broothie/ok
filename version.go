package ok

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed VERSION
var version string

func Version() string {
	return fmt.Sprintf("v%s", strings.TrimSpace(version))
}
