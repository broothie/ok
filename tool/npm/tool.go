package npm

import (
	"encoding/json"
	"os"

	"github.com/broothie/ok/task"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

type Tool struct{}

func (Tool) Name() string {
	return "NPM"
}

func (Tool) Executable() string {
	return "npm"
}

func (Tool) Filenames() []string {
	return []string{"package.json"}
}

func (Tool) Extensions() []string {
	return nil
}

func (Tool) ProcessFile(path string) ([]task.Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	var packageJSON packageJSON
	if err := json.NewDecoder(file).Decode(&packageJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse file")
	}

	return lo.Map(lo.Keys(packageJSON.Scripts), func(name string, _ int) task.Task { return Task{name: name} }), nil
}
