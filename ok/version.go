package ok

import (
	"fmt"
	"io"
)

const Version = "v0.1.8"

func WriteVersion(w io.Writer) {
	fmt.Fprintf(w, "👌 ok %s\n", Version)
}
