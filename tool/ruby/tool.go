package ruby

import (
	"os"

	"github.com/broothie/okay/tool"
)

const (
	ToolName = "ruby"
	filename = "Okayfile.rb"
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
