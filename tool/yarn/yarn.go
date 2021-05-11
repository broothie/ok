package yarn

import (
	"os"

	"github.com/broothie/okay/tool"
)

const (
	ToolName = "yarn"
	filename = "package.json"
)

type Yarn struct{}

func (Yarn) Init() error {
	_, err := os.Create(filename)
	return err
}

func (Yarn) Check() error {
	return tool.Check(ToolName)
}
