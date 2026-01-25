package ok

import (
	"context"

	"github.com/bobg/errors"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

func (o *Ok) SetUpTasks(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	for _, tl := range o.tools {
		for _, filePath := range tl.FilePaths {
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
		toolTasks, err := tl.ProcessFile(ctx, filePath)
		if err != nil {
			return errors.Wrapf(err, "parsing file %q for tool %q", filePath, tl.Name)
		}

		tasks := lo.Map(toolTasks, func(task Task, _ int) taskInfo {
			return taskInfo{
				Task:     task,
				Tool:     &tl,
				FilePath: filePath,
			}
		})

		o.tasksLock.Lock()
		o.Tasks = append(o.Tasks, tasks...)
		o.tasksLock.Unlock()
		return nil
	}
}
