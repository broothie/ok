package ruby

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
	filename   string
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Parameters() parameter.Parameters {
	return t.parameters
}

func (t Task) Run(ctx context.Context, args argument.Arguments) error {
	if err := util.CommandContext(ctx, "ruby", "-e", t.generatedRubyCode(args)).Run(); err != nil {
		return errors.Wrap(err, "failed to run ruby command")
	}

	return nil
}

func (t Task) generatedRubyCode(args argument.Arguments) string {
	var argStrings []string
	for _, arg := range args.Required() {
		switch arg.Type {
		case parameter.TypeBool, parameter.TypeInt, parameter.TypeFloat:
			argStrings = append(argStrings, arg.Value)
		case parameter.TypeString:
			argStrings = append(argStrings, fmt.Sprintf("%q", arg.Value))
		}
	}

	for _, arg := range args.Optional() {
		switch arg.Type {
		case parameter.TypeBool, parameter.TypeInt, parameter.TypeFloat:
			argStrings = append(argStrings, fmt.Sprintf("%s: %s", arg.Name, arg.Value))
		case parameter.TypeString:
			argStrings = append(argStrings, fmt.Sprintf("%s: %q", arg.Name, arg.Value))
		}
	}

	argString := strings.Join(argStrings, ", ")
	return fmt.Sprintf("at_exit { %s(%s) }\n\nrequire './%s'", t.name, argString, t.filename)
}
