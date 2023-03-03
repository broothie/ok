package dockercompose

import (
	"os"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

const initContents = `version: "3.8"

services: {}
`

type dockerCompose struct {
	Services map[string]any `yaml:"services"`
}

type Tool struct {
	config *tool.Config
}

func New() tool.Tool {
	return Tool{
		config: &tool.Config{
			"filenames":  "docker-compose.yml,docker-compose.yaml",
			"executable": "docker",
		},
	}
}

func (Tool) Name() string {
	return "docker-compose"
}

func (t Tool) Config() *tool.Config {
	return t.config
}

func (Tool) Init() error {
	return util.InitFile("docker-compose.yml", []byte(initContents))
}

func (t Tool) ProcessFile(path string) ([]task.Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	var dockerCompose dockerCompose
	if err := yaml.NewDecoder(file).Decode(&dockerCompose); err != nil {
		return nil, errors.Wrap(err, "failed to parse file")
	}

	return lo.Map(lo.Keys(dockerCompose.Services), func(service string, _ int) task.Task {
		return Task{
			Tool: t,
			name: service,
		}
	}), nil
}
