package make

import "os"

const (
	ToolName = "make"
	filename = "Makefile"
)

type Make struct{}

func (Make) Init() error {
	_, err := os.Create(filename)
	return err
}
