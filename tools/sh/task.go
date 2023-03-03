package sh

import (
	"context"
	"fmt"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Task struct {
	Tool
	name        string
	description string
	shCode      *string
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Description() string {
	return t.description
}

func (t Task) Parameters() task.Parameters {
	return task.NewSplatParameters(task.TypeString)
}

func (t Task) Run(ctx context.Context, args task.Arguments) error {
	if err := util.CommandContext(ctx, t.Config().Executable(), "-c", t.generatedShCode(args)).Run(); err != nil {
		return errors.Wrap(err, "failed to run sh command")
	}

	return nil
}

func (t Task) generatedShCode(args task.Arguments) string {
	argsString := strings.Join(lo.Map(args, func(arg task.Argument, _ int) string { return arg.Value }), " ")
	return fmt.Sprintf("%s\n\n%s %s", *t.shCode, t.name, argsString)
}
