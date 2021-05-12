package tool

import "fmt"

func List() {
	for _, tool := range Registry {
		if err := tool.Check(); err != nil {
			fmt.Printf("𝘹 %s %v\n", tool.Name(), err)
		} else {
			fmt.Printf("✔ %s\n", tool.Name())
		}
	}
}
