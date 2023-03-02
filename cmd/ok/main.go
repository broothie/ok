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
	"github.com/broothie/ok/cli"
	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/task"
	"github.com/radovskyb/watcher"
)

func main() {
	parser := cli.NewFromArgs()
	okArgs, err := parser.ParseOkArgs()
	if err != nil {
		logger.Log.Fatalf("failed to parse args: %v", err)
	}

	if okArgs.Help {
		if err := cli.Help(); err != nil {
			logger.Log.Fatalf("failed to print help: %v", err)
		}

		return
	}

	if okArgs.Version {
		fmt.Println(ok.Version())
		return
	}

	app := app.NewAsConfigured()
	if okArgs.ListTools {
		if err := app.ListTools(); err != nil {
			logger.Log.Fatalf("failed to list tools: %v", err)
		}

		return
	}

	if okArgs.TaskName == "" {
		if err := app.ListTasks(); err != nil {
			logger.Log.Fatalf("failed to list tasks: %v", err)
		}

		return
	}

	task, found := app.Task(okArgs.TaskName)
	if !found {
		logger.Log.Fatalf("no task found with name %q", okArgs.TaskName)
	}

	args, err := parser.ParseTaskArgs(task.Parameters())
	if err != nil {
		logger.Log.Fatalf("failed to parse task args: %v", err)
	}

	if len(okArgs.Watches) == 0 {
		if err := task.Run(context.Background(), args); err != nil {
			logger.Log.Printf("failed to run task %q: %v", task.Name(), err)
		}

		return
	}

	watcher := newWatcher(okArgs.Watches)
	defer watcher.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		taskCtx, taskCancel := context.WithCancel(ctx)
		go runTask(taskCtx, task, args)

		for range watcher.Event {
			taskCancel()
			taskCtx, taskCancel = context.WithCancel(ctx)
			go runTask(taskCtx, task, args)
		}
	}()

	if err := watcher.Start(100 * time.Millisecond); err != nil {
		logger.Log.Fatalf("failed to start watching files: %v", err)
	}
}

func runTask(ctx context.Context, task task.Task, args task.Arguments) {
	if err := task.Run(ctx, args); err != nil && !strings.HasSuffix(err.Error(), "signal: killed") {
		logger.Log.Printf("failed to run task %q: %v", task.Name(), err)
	}
}

func newWatcher(watches []string) *watcher.Watcher {
	watcher := watcher.New()
	watcher.SetMaxEvents(1)
	defer watcher.Close()

	for _, watch := range watches {
		matches, err := doublestar.FilepathGlob(watch)
		if err != nil {
			logger.Log.Fatalf("failed to glob pattern %q: %v", watch, err)
		}

		for _, match := range matches {
			stat, err := os.Stat(match)
			if err != nil {
				logger.Log.Fatalf("failed to stat file %q: %v", match, err)
			}

			if stat.IsDir() {
				if err := watcher.AddRecursive(match); err != nil {
					logger.Log.Fatalf("failed to add dir %q to watches: %v", match, err)
				}
			} else {
				if err := watcher.Add(match); err != nil {
					logger.Log.Fatalf("failed to add file %q to watches: %v", match, err)
				}
			}
		}
	}

	return watcher
}
