package tools

import (
	"fmt"

	"github.com/broothie/ok/tool"
)

func Index() {
	for _, tool := range tool.Registry {
		if err := tool.Check(); err != nil {
			fmt.Printf("𝘹 %s %v\n", tool.Name(), err)
		} else {
			fmt.Printf("✔ %s\n", tool.Name())
		}
	}
}
