package golang

import (
	"fmt"
	"os"

	"github.com/broothie/ok/util"
	"github.com/pelletier/go-toml"
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
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file %q already exists", filename)
		}

		return errors.Wrapf(err, "failed to create file %q", filename)
	}

	if _, err := fmt.Fprint(file, initContents); err != nil {
		return errors.Wrapf(err, "failed to init file %q", filename)
	}

	return errors.Wrapf(file.Close(), "failed to close file %q", filename)
}

func (Tool) Check() error {
	return util.Check(ToolName)
}

func (Tool) Configure(*toml.Decoder) error {
	return nil
}
