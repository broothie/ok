package npm

import (
	"encoding/json"
	"os"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

const initContents = `{
	"scripts": {}
}`

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

type Tool struct {
	config *tool.Config
}

func New() tool.Tool {
	return Tool{
		config: &tool.Config{
			"filenames":  "package.json",
			"executable": "npm",
		},
	}
}

func (Tool) Name() string {
	return "npm"
}

func (t Tool) Config() *tool.Config {
	return t.config
}

func (Tool) Init() error {
	return util.InitFile("package.json", []byte(initContents))
}

func (t Tool) ProcessFile(path string) ([]task.Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	var packageJSON packageJSON
	if err := json.NewDecoder(file).Decode(&packageJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse file")
	}

	return lo.Map(lo.Entries(packageJSON.Scripts), func(script lo.Entry[string, string], _ int) task.Task {
		return Task{
			Tool:        t,
			name:        script.Key,
			description: script.Value,
		}
	}), nil
}
