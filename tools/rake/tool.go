package rake

import (
	"fmt"
	"os"

	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

const (
	ToolName = "rake"
	filename = "Rakefile"
)

type Tool struct {
	ToolConfig Config
}

type Config struct {
	Bundle *bool `toml:"bundle"`
}

func (Tool) Name() string {
	return ToolName
}

func (Tool) Init() error {
	_, err := os.OpenFile(filename, os.O_CREATE, 0666)
	if os.IsExist(err) {
		return fmt.Errorf("file '%s' already exists", filename)
	}

	return errors.Wrapf(err, "could not create file '%s'", filename)
}

func (Tool) Check() error {
	return util.Check(ToolName)
}

func (t *Tool) Config() interface{} {
	return &t.ToolConfig
}
