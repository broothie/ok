package ok

import (
	"context"
	"os"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/samber/lo"
)

func (o *Ok) RunTask(ctx context.Context, taskName string, remainingArgs []string) error {
	tsk, found := lo.Find(o.Tasks, func(tsk taskInfo) bool { return tsk.Name == taskName })
	if !found {
		return errors.Errorf("no task found with name %q", taskName)
	}

	return errors.Wrapf(tsk.Run(ctx, remainingArgs), "running task %q", taskName)
}

func (i taskInfo) Run(ctx context.Context, remainingArgs []string) error {
	options, err := i.RunOptions(ctx, remainingArgs)
	if err != nil {
		return errors.Wrapf(err, "getting run options for %s task %q", i.Tool.Name, i.Name)
	}

	options = append(options,
		cob.SetStdin(os.Stdin),
		cob.SetStdout(os.Stdout),
		cob.SetStderr(os.Stderr),
	)

	_, err = cob.Run(ctx, i.Tool.CommandPath, options...)
	return errors.Wrapf(err, "running %s task %q", i.Tool.Name, i.Name)
}
