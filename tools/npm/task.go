package npm

import (
	"context"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Task struct {
	Tool
	description string
	name        string
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Description() string {
	return t.description
}

func (t Task) Parameters() task.Parameters {
	return nil
}

func (t Task) Run(ctx context.Context, args task.Arguments) error {
	commandArgs := []string{"run", t.name}
	commandArgs = append(commandArgs, lo.Map(args, func(arg task.Argument, _ int) string { return arg.Value })...)
	if err := util.CommandContext(ctx, t.Config().Executable(), commandArgs...).Run(); err != nil {
		return errors.Wrap(err, "failed to run npm script")
	}

	return nil
}
