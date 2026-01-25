package application

import (
	"context"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/bobg/errors"
	"github.com/broothie/ok/tool"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

func (a *Application) DiscoverTasks(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	for _, tl := range a.Tools {
		if err := ctx.Err(); err != nil {
			return err
		}

		group.Go(a.globFiles(ctx, tl, group))
	}

	return group.Wait()
}

func (a *Application) globFiles(ctx context.Context, tl tool.Tool, group *errgroup.Group) func() error {
	return func() error {
		seen := make(map[string]bool)

		for _, glob := range tl.FileGlobs {
			if err := ctx.Err(); err != nil {
				return err
			}

			filePaths, err := doublestar.Glob(os.DirFS("."), glob)
			if err != nil {
				return errors.Wrapf(err, "globbing %q for tool %q", glob, tl.Name)
			}

			for _, filePath := range filePaths {
				if err := ctx.Err(); err != nil {
					return err
				}

				// Dedupe by lowercase path for case-insensitive filesystems
				key := strings.ToLower(filePath)
				if seen[key] {
					continue
				}
				seen[key] = true

				group.Go(a.parseFile(ctx, tl, filePath))
			}
		}

		return nil
	}
}

func (a *Application) parseFile(ctx context.Context, tl tool.Tool, filePath string) func() error {
	return func() error {
		toolTasks, err := tl.ParseFile(ctx, filePath)
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
