package ruby

import (
	"os"

	"github.com/broothie/okay/tool"
)

const (
	ToolName = "ruby"
	filename = "Okayfile.rb"
)

type Ruby struct{}

func (Ruby) Init() error {
	_, err := os.Create(filename)
	return err
}

func (Ruby) Check() error {
	return tool.Check(ToolName)
}
