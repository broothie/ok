package node

import (
	"context"
	"fmt"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

type Task struct {
	Tool
	name        string
	description string
	parameters  task.Parameters
	jsCode      *string
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Description() string {
	return t.description
}

func (t Task) Parameters() task.Parameters {
	return t.parameters
}

func (t Task) Run(ctx context.Context, args task.Arguments) error {
	if err := util.CommandContext(ctx, t.Config().Executable(), "-e", t.generatedJSCode(args)).Run(); err != nil {
		return errors.Wrap(err, "failed to run node command")
	}

	return nil
}

func (t Task) generatedJSCode(args task.Arguments) string {
	var argStrings []string
	for _, arg := range args.Required() {
		switch arg.Type {
		case task.TypeBool, task.TypeInt, task.TypeFloat:
			argStrings = append(argStrings, arg.Value)
		case task.TypeString:
			argStrings = append(argStrings, fmt.Sprintf("%q", arg.Value))
		}
	}

	for _, arg := range args.Optional() {
		switch arg.Type {
		case task.TypeBool, task.TypeInt, task.TypeFloat:
			argStrings = append(argStrings, fmt.Sprintf("%s=%s", arg.Name, arg.Value))
		case task.TypeString:
			argStrings = append(argStrings, fmt.Sprintf("%s=%q", arg.Name, arg.Value))
		}
	}

	argString := strings.Join(argStrings, ", ")
	return fmt.Sprintf("%s\n\n%s(%s)", *t.jsCode, t.name, argString)
}
