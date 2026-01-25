package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bobg/errors"
	"github.com/broothie/ok"
	"github.com/broothie/ok/cli"
	"github.com/broothie/ok/tool/tools"
	"github.com/samber/lo"
)

const version = "v0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func run() error {
	helpFlag := &cli.Flag{
		Name:         "help",
		Type:         cli.FlagTypeBool,
		Help:         "Show command help.",
		Shorts:       []rune{'h'},
		DefaultValue: false,
	}

	directoryFlag := &cli.Flag{
		Name:         "directory",
		Type:         cli.FlagTypeString,
		Help:         "Directory to run command from.",
		Shorts:       []rune{'d'},
		DefaultValue: ".",
	}

	timeoutFlag := &cli.Flag{
		Name:         "timeout",
		Type:         cli.FlagTypeDuration,
		Help:         "Command timeout.",
		DefaultValue: time.Second,
	}

	toolFlag := &cli.Flag{
		Name:         "tools",
		Type:         cli.FlagTypeBool,
		Help:         "List tools.",
		DefaultValue: false,
	}

	parser := cli.NewParser(version, os.Args[1:], helpFlag, directoryFlag, timeoutFlag, toolFlag)
	if err := parser.Parse(); err != nil {
		return errors.Wrap(err, "parsing flags")
	}

	if err := os.Chdir(directoryFlag.Value().(string)); err != nil {
		return errors.Wrapf(err, "changing dir to %q", directoryFlag.Value().(string))
	}

	if helpFlag.Value().(bool) {
		return errors.Wrap(parser.WriteHelp(os.Stdout), "printing help")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutFlag.Value().(time.Duration))
	defer cancel()

	app := ok.New()

	if err := app.SetUpTools(ctx, tools.All()); err != nil {
		return errors.Wrap(err, "setting up tools")
	}

	if toolFlag.Value().(bool) {
		return errors.Wrap(app.PrintTools(os.Stdout), "listing tools")
	}

	if err := app.DiscoverTasks(ctx); err != nil {
		return errors.Wrap(err, "discovering tasks")
	}

	taskName := parser.TaskName()
	if taskName == "" {
		return errors.Wrap(app.ListTasks(os.Stdout), "listing tasks")
	}

	tsk, found := lo.Find(app.Tasks, func(task ok.Task) bool { return task.Name == taskName })
	if !found {
		return errors.Errorf("no task found with name %q", taskName)
	}

	if err := tsk.Run(ctx, parser.RemainingArgs()); err != nil {
		return errors.Wrapf(err, "running task %q", tsk.Name)
	}

	return nil
}
