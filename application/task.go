package application

import (
	"context"
	"os"
	"os/exec"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok/tool"
)

type Task struct {
	tool.Task
	Tool     *tool.Tool
	FilePath string
}

func (t Task) Run(ctx context.Context, args []string) error {
	commandPath, err := exec.LookPath(t.Tool.CommandName)
	if err != nil {
		return errors.Wrapf(err, "looking up command %q", t.Tool.CommandName)
	}

	options, err := t.RunOptions(ctx, args)
	if err != nil {
		return errors.Wrapf(err, "getting run options for %s task %q", t.Tool.Name, t.Name)
	}

	options = append(options,
		cob.SetStdin(os.Stdin),
		cob.SetStdout(os.Stdout),
		cob.SetStderr(os.Stderr),
	)

	_, err = cob.Run(ctx, commandPath, options...)
	return errors.Wrapf(err, "running %s task %q", t.Tool.Name, t.Name)
}
