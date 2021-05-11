package make

import (
	"os"

	"github.com/broothie/okay/tool"
)

const (
	ToolName = "make"
	filename = "Makefile"
)

type Make struct{}

func (Make) Init() error {
	_, err := os.Create(filename)
	return err
}

func (Make) Check() error {
	return tool.Check(ToolName)
}
