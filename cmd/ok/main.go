package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/broothie/ok"
	"github.com/broothie/ok/app"
	"github.com/broothie/ok/arg"
	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/task"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	app := app.New(app.Tools())
	parser := arg.NewParser(args)

	// Parse options
	options, err := app.Options(parser)
	if err != nil {
		return err
	}

	// Handle early exits
	if options.Help {
		return app.PrintHelp()
	} else if options.Version {
		fmt.Println(ok.Version())
		return nil
	} else if options.ListTools {
		return app.Tools.Print()
	}

	// List task when no task provided
	if options.TaskName == "" {
		return app.Tasks().Print()
	}

	task, found := app.Tasks().Task(options.TaskName)
	if !found {
		return fmt.Errorf("unknown task %q", options.TaskName)
	}

	// Parse task args
	taskArgs, err := task.Parameters().Parse(parser)
	if err != nil {
		return err
	}

	// Run task in foreground if no watches provided
	if len(options.Watches) == 0 {
		return task.Run(context.Background(), taskArgs)
	}

	// Set up file watcher
	fileWatcher, err := newWatcher(options.Watches)
	if err != nil {
		return err
	}
	defer fileWatcher.Close()

	// Run task in background, killing/restarting on each watch event
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		taskCtx, taskCancel := context.WithCancel(ctx)
		go backgroundTask(taskCtx, task, taskArgs)

		for range fileWatcher.Event {
			taskCancel()
			taskCtx, taskCancel = context.WithCancel(ctx)
			go backgroundTask(taskCtx, task, taskArgs)
		}
	}()

	// Block foreground on file watching
	if err := fileWatcher.Start(100 * time.Millisecond); err != nil {
		return errors.Wrap(err, "failed to start watching files")
	}

	return nil
}

func backgroundTask(ctx context.Context, task task.Task, args task.Arguments) {
	if err := task.Run(ctx, args); err != nil && !strings.HasSuffix(err.Error(), "signal: killed") {
		logger.Log.Printf("failed to run task %q: %v", task.Name(), err)
	}
}

func newWatcher(watches []string) (*watcher.Watcher, error) {
	watcher := watcher.New()
	watcher.SetMaxEvents(1)

	for _, watch := range watches {
		matches, err := doublestar.FilepathGlob(watch)
		if err != nil {
			return nil, errors.Wrap(err, "failed to glob pattern")
		}

		for _, match := range matches {
			stat, err := os.Stat(match)
			if err != nil {
				return nil, errors.Wrap(err, "failed to state file")
			}

			if stat.IsDir() {
				if err := watcher.AddRecursive(match); err != nil {
					return nil, errors.Wrap(err, "failed to add dir to watches")
				}
			} else {
				if err := watcher.Add(match); err != nil {
					return nil, errors.Wrap(err, "failed to add file to watches")
				}
			}
		}
	}

	return watcher, nil
}
