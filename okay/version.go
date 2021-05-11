package okay

import (
	"fmt"
	"io"
)

const Version = "0.1.0"

func WriteVersion(w io.Writer) {
	fmt.Fprintf(w, "ok v%s 👌\n", Version)
}
