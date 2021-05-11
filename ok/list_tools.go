package ok

import (
	"fmt"
	"io"
)

func ListTools(w io.Writer) {
	for _, tool := range Registry {
		if err := tool.Check(); err != nil {
			fmt.Fprintf(w, "𝘹 %s %v\n", tool.Name(), err)
		} else {
			fmt.Fprintf(w, "✔ %s\n", tool.Name())
		}
	}
}
