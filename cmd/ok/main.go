package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/broothie/ok"
	clipkg "github.com/broothie/ok/cli"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tools"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	okrcArgs, err := okrc()
	if err != nil {
		return err
	}

	// Parse options
	cli := clipkg.New(append(okrcArgs, os.Args[1:]...))
	options, err := cli.Options()
	if err != nil {
		return err
	}

	tools := tools.FromRegistry()

	// Handle tool config
	for _, options := range options.Tool {
		if options.Action() == clipkg.ToolOptionsActionSet {
			tool, toolPresent := tools[options.Name]
			if !toolPresent {
				return fmt.Errorf("unknown tool %q", options.Name)
			}

			tool.Config().Set(options.Key, options.Value)
		}
	}

	// Handle early exit tool options
	if shown, err := handleShowToolConfig(tools, options); err != nil {
		return err
	} else if shown {
		return nil
	}

	// Handle early exits
	if options.Help {
		return cli.PrintHelp()
	} else if options.Version {
		fmt.Println(ok.Version())
		return nil
	} else if options.ListTools {
		return tools.Print()
	}

	// Fetch tasks
	tasks, err := tools.CollectTasks()
	if err != nil {
		return err
	}

	// List task when no task provided
	if options.TaskName == "" {
		return tasks.Print()
	}

	task, found := tasks.Task(options.TaskName)
	if !found {
		return fmt.Errorf("unknown task %q", options.TaskName)
	}

	// Parse task args
	args, err := cli.ParseParameters(task.Parameters())
	if err != nil {
		return err
	}

	// Run task in foreground if no watches provided
	if len(options.Watches) == 0 {
		return task.Run(context.Background(), args)
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
		go backgroundTask(taskCtx, task, args)

		for range fileWatcher.Event {
			taskCancel()
			taskCtx, taskCancel = context.WithCancel(ctx)
			go backgroundTask(taskCtx, task, args)
		}
	}()

	// Block foreground on file watching
	if err := fileWatcher.Start(100 * time.Millisecond); err != nil {
		return errors.Wrap(err, "failed to start watching files")
	}

	return nil
}

func okrc() ([]string, error) {
	contents, err := os.ReadFile(".okrc")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to read .okrc")
	}

	return strings.Fields(string(contents)), nil
}

func handleShowToolConfig(tools tools.Tools, options clipkg.Options) (bool, error) {
	for _, toolOptions := range options.Tool {
		switch toolOptions.Action() {
		case clipkg.ToolOptionsActionTools:
			for _, tool := range tools {
				for key, value := range tool.Config().Entries() {
					fmt.Printf("--tool %s.%s=%s\n", tool.Name(), key, value)
				}
			}

			return true, nil

		case clipkg.ToolOptionsActionTool:
			tool, toolPresent := tools[toolOptions.Name]
			if !toolPresent {
				return false, fmt.Errorf("unknown tool %q", toolOptions.Name)
			}

			for key, value := range tool.Config().Entries() {
				fmt.Printf("--tool %s.%s=%s\n", tool.Name(), key, value)
			}

			return true, nil

		case clipkg.ToolOptionsActionKey:
			tool, toolPresent := tools[toolOptions.Name]
			if !toolPresent {
				return false, fmt.Errorf("unknown tool %q", toolOptions.Name)
			}

			value := tool.Config().Get(toolOptions.Key)
			fmt.Printf("--tool %s.%s=%s\n", tool.Name(), toolOptions.Key, value)
			return true, nil
		}
	}

	return false, nil
}

func backgroundTask(ctx context.Context, task task.Task, args task.Arguments) {
	if err := task.Run(ctx, args); err != nil && !strings.HasSuffix(err.Error(), "signal: killed") {
		fmt.Println("failed to run task", err)
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
