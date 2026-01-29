package ok

import (
	"context"
	"os"
	"sync"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/samber/lo"
)

type Ok struct {
	toolsLock *sync.Mutex
	tools     []toolInfo
	tasksLock *sync.Mutex
	tasks     []taskInfo
}

func New() *Ok {
	return &Ok{
		toolsLock: new(sync.Mutex),
		tasksLock: new(sync.Mutex),
	}
}

func (o *Ok) RunTask(ctx context.Context, taskName string, remainingArgs []string) error {
	tsk, found := lo.Find(o.tasks, func(tsk taskInfo) bool { return tsk.Name == taskName })
	if !found {
		return errors.Errorf("no task found with name %q", taskName)
	}

	return tsk.Run(ctx, remainingArgs)
}

func (i taskInfo) Run(ctx context.Context, remainingArgs []string) error {
	toolCfg := ToolConfig{
		Executable: i.tool.commandPath,
	}
	options, err := i.RunOptions(ctx, remainingArgs, toolCfg)
	if err != nil {
		return errors.Wrapf(err, "getting run options for %s task %q", i.tool.Name, i.Name)
	}

	options = append(options,
		cob.AddEnv("OK_TASK", i.Task.Name),
		cob.AddEnv("OK_TOOL", i.tool.Name),
		cob.AddEnv("OK_VERSION", Version()),
		cob.SetStdin(os.Stdin),
		cob.SetStdout(os.Stdout),
		cob.SetStderr(os.Stderr),
	)

	_, err = cob.Run(ctx, i.tool.commandPath, options...)
	return errors.Wrapf(err, "running %s task %q", i.tool.Name, i.Name)
}
