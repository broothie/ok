package rake

import (
	"fmt"
	"os"

	"github.com/broothie/ok/util"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

const (
	ToolName = "rake"
	filename = "Rakefile"
)

type Tool struct {
	Config Config
}

type Config struct {
	Bundler *bool `toml:"bundler"`
}

func (Tool) Name() string {
	return ToolName
}

func (Tool) Init() error {
	_, err := os.OpenFile(filename, os.O_CREATE, 0666)
	if os.IsExist(err) {
		return fmt.Errorf("file %q already exists", filename)
	}

	return errors.Wrapf(err, "could not create file %q", filename)
}

func (Tool) Check() error {
	return util.Check(ToolName)
}

func (t *Tool) Configure(decoder *toml.Decoder) error {
	return decoder.Decode(&t.Config)
}
