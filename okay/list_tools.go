package okay

import (
	"fmt"
	"io"
)

func ListTools(w io.Writer) {
	for toolName, tool := range Registry {
		if err := tool.Check(); err != nil {
			fmt.Fprintf(w, "𝘹 %s %v\n", toolName, err)
		} else {
			fmt.Fprintf(w, "✔ %s\n", toolName)
		}
	}
}
