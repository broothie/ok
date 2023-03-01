package npm

import (
	"context"

	"github.com/broothie/ok/argument"
	"github.com/broothie/ok/parameter"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

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
	if err := util.CommandContext(ctx, "npm", commandArgs...).Run(); err != nil {
		return errors.Wrap(err, "failed to run npm script")
	}

	return nil
}
