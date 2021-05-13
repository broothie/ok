package ruby

import (
	"fmt"
	"os"

	"github.com/broothie/ok/toolhelp"
	"github.com/pkg/errors"
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
	_, err := os.OpenFile(filename, os.O_CREATE, 0)
	if os.IsExist(err) {
		return fmt.Errorf("file '%s' already exists", filename)
	}

	return errors.Wrapf(err, "could not create file '%s'", filename)
}

func (Tool) Check() error {
	return toolhelp.Check(ToolName)
}
