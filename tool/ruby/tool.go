package ruby

import (
	"os"

	"github.com/broothie/ok/tool"
)

const (
	ToolName = "ruby"
	filename = "Okfile.rb"
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
