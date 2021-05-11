package dockercompose

import (
	"fmt"
	"os"

	"github.com/broothie/okay/tool"
	"github.com/pkg/errors"
)

const (
	ToolName = "docker-compose"
	filename = "docker-compose.yml"

	initContents = `version: "3"
services:
`
)

type Tool struct{}

func (Tool) Name() string {
	return ToolName
}

func (Tool) Init() error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "failed to create file '%s'", filename)
	}

	if _, err := fmt.Fprint(file, initContents); err != nil {
		return errors.Wrapf(err, "failed to init file '%s'", filename)
	}

	return errors.Wrapf(file.Close(), "failed to close file '%s'", filename)
}

func (Tool) Check() error {
	return tool.Check(ToolName)
}
