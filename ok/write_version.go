package ok

import (
	"fmt"
	"io"
)

func WriteVersion(w io.Writer) {
	fmt.Fprintf(w, "👌 ok %s\n", Version)
}
