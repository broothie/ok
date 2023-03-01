package python

import (
	"context"
	"fmt"
	"strings"

	"github.com/broothie/ok/argument"
	"github.com/broothie/ok/parameter"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

type Task struct {
	name       string
	parameters parameter.Parameters
	pythonCode *string
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Parameters() parameter.Parameters {
	return t.parameters
}

func (t Task) Run(ctx context.Context, args argument.Arguments) error {
	if err := util.CommandContext(ctx, "python3", "-c", t.generatedPythonCode(args)).Run(); err != nil {
		return errors.Wrap(err, "failed to run python command")
	}

	return nil
}

func (t Task) generatedPythonCode(args argument.Arguments) string {
	var argStrings []string
	for _, arg := range args.Required() {
		switch arg.Type {
		case parameter.TypeBool:
			argStrings = append(argStrings, fmt.Sprintf("%s%s", strings.ToTitle(string(arg.Value[0])), arg.Value[1:]))
		case parameter.TypeInt, parameter.TypeFloat:
			argStrings = append(argStrings, arg.Value)
		case parameter.TypeString:
			argStrings = append(argStrings, fmt.Sprintf("%q", arg.Value))
		}
	}

	for _, arg := range args.Optional() {
		switch arg.Type {
		case parameter.TypeBool:
			argStrings = append(argStrings, fmt.Sprintf("%s=%s%s", arg.Name, strings.ToTitle(string(arg.Value[0])), arg.Value[1:]))
		case parameter.TypeInt, parameter.TypeFloat:
			argStrings = append(argStrings, fmt.Sprintf("%s=%s", arg.Name, arg.Value))
		case parameter.TypeString:
			argStrings = append(argStrings, fmt.Sprintf("%s=%q", arg.Name, arg.Value))
		}
	}

	argString := strings.Join(argStrings, ", ")
	return fmt.Sprintf("%s\n\n%s(%s)", *t.pythonCode, t.name, argString)
}
