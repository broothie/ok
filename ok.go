package main

import (
	_ "embed"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/broothie/ok/cli"
	"github.com/broothie/ok/ok"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tools"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
)

//go:embed VERSION
var version string

func Version() string {
	return strings.TrimSpace(version)
}

func main() {
	// Parse options
	parser, err := cli.NewParser(os.Args[1:])
	if err != nil {
		ok.Logger.Println(err)
		os.Exit(1)
		return
	}

	options, err := parser.ParseOptions()
	if err != nil {
		ok.Logger.Println(err)
		os.Exit(1)
		return
	}

	if options.Debug {
		ok.DebugLogger.Printf("%+v\n", options)
	}

	// Process options
	switch {
	case options.Help:
		if err := cli.PrintHelp(Version()); err != nil {
			ok.Logger.Println(err)
			os.Exit(1)
			return
		}

	case options.Version:
		cli.PrintVersion(Version())

	case options.Init != "":
		if err := tools.InitTool(options.Init); err != nil {
			ok.Logger.Println(err)
			os.Exit(1)
			return
		}

	case options.ListTools:
		tools.List()

	case options.TaskName == "":
		if err := task.List(tools.Mount(options.ConfigurableOptions.SkipTools(), options.ConfigurableOptions.ToolSort())); err != nil {
			ok.Logger.Println(err)
			os.Exit(1)
			return
		}
	}

	if options.Stop {
		return
	}

	// Get task
	taskName := options.TaskName
	tasks := tools.Mount(options.ConfigurableOptions.SkipTools(), options.ConfigurableOptions.ToolSort())
	task, taskExists := tasks[options.TaskName]
	if !taskExists {
		ok.Logger.Printf("no task called '%s'", taskName)
		os.Exit(1)
		return
	}

	// Parse task args
	args, err := parser.ParseArgs(task.Params())
	if err != nil {
		ok.Logger.Println(err)
		os.Exit(1)
		return
	}

	// Run task
	if len(options.Watches) > 0 {
		if err := runWatcher(task, args, options.Watches); err != nil {
			ok.Logger.Println(err)
			os.Exit(1)
		}
	} else {
		if err := task.Invoke(args).Wait(); err != nil {
			if err.Error() == "wait: no child processes" {
				return
			}

			ok.Logger.Println(err)
			os.Exit(1)
		}
	}
}

func runWatcher(t task.Task, args task.Args, watches []string) error {
	var wg sync.WaitGroup
	defer wg.Wait()
	watcher := watcher.New()
	watcher.SetMaxEvents(1)

	for _, watchPattern := range watches {
		filenames, err := doublestar.Glob(watchPattern)
		if err != nil {
			return errors.Wrapf(err, "failed to glob '%s'", watchPattern)
		}

		for _, filename := range filenames {
			if err := watcher.Add(filename); err != nil {
				return errors.Wrapf(err, "failed to add file '%s' to watches", filename)
			}
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var process task.RunningTask

		for {
			select {
			case <-watcher.Event:
				if process != nil {
					process.Kill()
				}

				process = t.Invoke(args)

			case err := <-watcher.Error:
				if process != nil {
					process.Kill()
				}

				ok.Logger.Println(err)

			case <-watcher.Closed:
				if process != nil {
					process.Kill()
				}

				return
			}
		}
	}()

	kill := make(chan os.Signal)
	signal.Notify(kill, os.Interrupt)
	signal.Notify(kill, os.Kill)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-kill
		watcher.Close()
	}()

	if err := watcher.Start(100 * time.Millisecond); err != nil {
		return errors.Wrap(err, "failed to start watcher")
	}

	return nil
}
