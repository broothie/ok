package make

import (
	"context"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

type Task struct {
	Tool
	name string
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Parameters() task.Parameters {
	return nil
}

func (t Task) Run(ctx context.Context, _ task.Arguments) error {
	if err := util.CommandContext(ctx, t.Config().Executable(), t.name).Run(); err != nil {
		return errors.Wrap(err, "failed to run make rule")
	}

	return nil
}
