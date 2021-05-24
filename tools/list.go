package tools

import (
	"fmt"

	tool2 "github.com/broothie/ok/tool"
)

func List() {
	for _, tool := range tool2.Registry {
		if err := tool.Check(); err != nil {
			fmt.Printf("𝘹 %s %v\n", tool.Name(), err)
		} else {
			fmt.Printf("✔ %s\n", tool.Name())
		}
	}
}
