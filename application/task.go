package application

import (
	"context"
	"os"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok/tool"
)

type Task struct {
	tool.Task
	Tool     *Tool
	FilePath string
}

func (t Task) Run(ctx context.Context, remainingArgs []string) error {
	options, err := t.RunOptions(ctx, remainingArgs)
	if err != nil {
		return errors.Wrapf(err, "getting run options for %s task %q", t.Tool.Name, t.Name)
	}

	options = append(options,
		cob.SetStdin(os.Stdin),
		cob.SetStdout(os.Stdout),
		cob.SetStderr(os.Stderr),
	)

	_, err = cob.Run(ctx, t.Tool.CommandPath, options...)
	return errors.Wrapf(err, "running %s task %q", t.Tool.Name, t.Name)
}
