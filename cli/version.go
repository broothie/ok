package cli

import (
	_ "embed"
	"fmt"
	"io"
)

func PrintVersion(w io.Writer, version string) error {
	_, err := fmt.Fprintf(w, "👌 ok v%s\n", version)
	return err
}
