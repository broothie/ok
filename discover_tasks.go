package ok

import (
	"context"

	"github.com/bobg/errors"
	"github.com/broothie/ok/tool"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

func (a *Application) DiscoverTasks(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	for _, tl := range a.tools {
		for _, filePath := range tl.FilePaths {
			if err := ctx.Err(); err != nil {
				return err
			}

			group.Go(a.processFile(ctx, tl, filePath))
		}

	}

	return group.Wait()
}

func (a *Application) processFile(ctx context.Context, tl Tool, filePath string) func() error {
	return func() error {
		toolTasks, err := tl.ProcessFile(ctx, filePath)
		if err != nil {
			return errors.Wrapf(err, "parsing file %q for tool %q", filePath, tl.Name)
		}

		tasks := lo.Map(toolTasks, func(task tool.Task, _ int) Task {
			return Task{
				Task:     task,
				Tool:     &tl,
				FilePath: filePath,
			}
		})

		a.tasksLock.Lock()
		a.Tasks = append(a.Tasks, tasks...)
		a.tasksLock.Unlock()
		return nil
	}
}
