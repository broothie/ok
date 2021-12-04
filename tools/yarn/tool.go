package yarn

import (
	"fmt"
	"os"

	"github.com/broothie/ok/util"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
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
	_, err := os.OpenFile(filename, os.O_CREATE, 0)
	if os.IsExist(err) {
		return fmt.Errorf("file %q already exists", filename)
	}

	return errors.Wrapf(err, "could not create file %q", filename)
}

func (Tool) Check() error {
	return util.Check(ToolName)
}

func (Tool) Configure(*toml.Decoder) error {
	return nil
}
