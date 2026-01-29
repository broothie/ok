package ok

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"sort"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok/table"
	"github.com/broothie/option"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type Task struct {
	Name       string
	RunOptions func(ctx context.Context, args []string, toolCfg ToolConfig) (option.Options[*exec.Cmd], error)
}

type taskInfo struct {
	Task
	tool     *toolInfo
	filePath string
}

func (i taskInfo) Run(ctx context.Context, remainingArgs []string) error {
	options, err := i.RunOptions(ctx, remainingArgs, ToolConfig{
		Executable: i.tool.commandPath,
	})
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

func (o *Ok) ListTasks(w io.Writer) error {
	rows := lo.Map(o.tasks, func(tsk taskInfo, _ int) []string { return []string{tsk.Name, tsk.tool.Name, tsk.filePath} })
	sort.Slice(rows, func(i, j int) bool { return rows[i][0] < rows[j][0] })

	rows = append([][]string{{"TASK", "TOOL", "FILE"}}, rows...)
	return errors.Wrap(table.Write(w, rows), "writing table")
}

func (o *Ok) SetUpTasks(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	for _, tl := range o.tools {
		if err := tl.err(); err != nil {
			slog.Debug("skipping task setup for errored tool", slog.String("tool", tl.Name), slog.Any("error", err))
			continue
		}

		for _, filePath := range tl.filePaths {
			if err := ctx.Err(); err != nil {
				return err
			}

			group.Go(o.processFile(ctx, tl, filePath))
		}

	}

	return group.Wait()
}

func (o *Ok) processFile(ctx context.Context, tl toolInfo, filePath string) func() error {
	return func() error {
		toolTasks, err := tl.ProcessFile(ctx, filePath, ToolConfig{
			Executable: tl.commandPath,
		})
		if err != nil {
			return errors.Wrapf(err, "parsing file %q for tool %q", filePath, tl.Name)
		}

		tasks := lo.Map(toolTasks, func(task Task, _ int) taskInfo {
			return taskInfo{
				Task:     task,
				tool:     &tl,
				filePath: filePath,
			}
		})

		o.tasksLock.Lock()
		o.tasks = append(o.tasks, tasks...)
		o.tasksLock.Unlock()
		return nil
	}
}

func (o *Ok) RunTask(ctx context.Context, taskName string, remainingArgs []string) error {
	tsk, found := lo.Find(o.tasks, func(tsk taskInfo) bool { return tsk.Name == taskName })
	if !found {
		return errors.Errorf("no task found with name %q", taskName)
	}

	return tsk.Run(ctx, remainingArgs)
}
