package npm

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/broothie/ok/argument"
	"github.com/broothie/ok/parameter"
	"github.com/broothie/ok/tool"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

type Tool struct{}

func (Tool) Name() string {
	return "npm"
}

func (Tool) CommandName() string {
	return "npm"
}

func (Tool) Filenames() []string {
	return []string{"package.json"}
}

func (Tool) Extensions() []string {
	return nil
}

func (Tool) ProcessFile(path string) ([]tool.Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	var packageJSON packageJSON
	if err := json.NewDecoder(file).Decode(&packageJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse file")
	}

	return lo.Map(lo.Keys(packageJSON.Scripts), func(name string, _ int) tool.Task { return Task{name: name} }), nil
}

type Task struct {
	name string
}

func (r Task) Name() string {
	return r.name
}

func (r Task) Parameters() parameter.Parameters {
	return nil
}

func (r Task) Run(ctx context.Context, args argument.Arguments) error {
	commandArgs := []string{"run", r.name}
	commandArgs = append(commandArgs, lo.Map(args, func(arg argument.Argument, _ int) string { return arg.Value })...)

	cmd := exec.CommandContext(ctx, "npm", commandArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to run npm script")
	}

	return nil
}
