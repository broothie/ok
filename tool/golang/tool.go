package golang

import (
	"fmt"
	"os"

	"github.com/broothie/ok/tool"
	"github.com/pkg/errors"
)

const (
	ToolName = "go"
	filename = "Okfile.go"

	initContents = `//+build ok

package main
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
