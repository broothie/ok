package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/broothie/ok"
	"github.com/broothie/ok/argument"
	"github.com/broothie/ok/cli"
	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/tool"
	"github.com/radovskyb/watcher"
)

func main() {
	parser := cli.NewFromArgs()
	okArgs, err := parser.ParseOkArgs()
	if err != nil {
		logger.Log.Fatalf("failed to parse args: %v", err)
	}

	if okArgs.Help {
		fmt.Println("helping...")
		return
	}

	ok := ok.NewAsConfigured()
	if okArgs.ListTools {
		if err := ok.ListTools(); err != nil {
			logger.Log.Fatalf("failed to list tools: %v", err)
		}

		return
	}

	if okArgs.TaskName == "" {
		if err := ok.ListTasks(); err != nil {
			logger.Log.Fatalf("failed to list tasks: %v", err)
		}

		return
	}

	task, found := ok.Task(okArgs.TaskName)
	if !found {
		logger.Log.Fatalf("no task found with name %q", okArgs.TaskName)
	}

	args, err := parser.ParseTaskArgs(task.Parameters())
	if err != nil {
		logger.Log.Fatalf("failed to parse task args: %v", err)
	}

	watcher := watcher.New()
	watcher.SetMaxEvents(1)
	defer watcher.Close()

	for _, watch := range okArgs.Watches {
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	go func() {
		<-interruptChan
		cancel()
	}()

	go func() {
		taskCtx, taskCancel := context.WithCancel(ctx)
		go runTask(taskCtx, task, args, cancel)

		for {
			select {
			case <-ctx.Done():
				watcher.Close()

			case <-watcher.Event:
				taskCancel()
				taskCtx, taskCancel = context.WithCancel(ctx)
				go runTask(taskCtx, task, args, cancel)
			}
		}
	}()

	if err := watcher.Start(100 * time.Millisecond); err != nil {
		logger.Log.Fatalf("failed to start watching files: %v", err)
	}
}

func runTask(ctx context.Context, task tool.Task, args argument.Arguments, done func()) {
	if err := task.Run(ctx, args); err != nil {
		if !strings.HasSuffix(err.Error(), "signal: killed") {
			logger.Log.Printf("failed to run task %q: %v", task.Name(), err)
		}

		return
	}

	done()
}
