package yarn

import (
	"os"

	"github.com/broothie/ok/tool"
)

const (
	ToolName = "yarn"
	filename = "package.json"
)

type Tool struct{}

func (Tool) Name() string {
	return ToolName
}

func (Tool) Init() error {
	_, err := os.Create(filename)
	return err
}

func (Tool) Check() error {
	return tool.Check(ToolName)
}
