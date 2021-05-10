package yarn

import "os"

const (
	ToolName = "yarn"
	filename = "package.json"
)

type Yarn struct{}

func (Yarn) Init() error {
	_, err := os.Create(filename)
	return err
}
