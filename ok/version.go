package ok

import (
	_ "embed"
	"fmt"
)

func PrintVersion(version string) {
	fmt.Printf("👌 ok v%s\n", version)
}
