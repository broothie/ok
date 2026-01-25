package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bobg/errors"
	"github.com/broothie/ok/application"
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
		Shorts:       []rune{'t'},
		DefaultValue: time.Second,
	}

	parser := cli.NewParser(version, os.Args[1:], helpFlag, directoryFlag, timeoutFlag)
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

	app := application.New(directoryFlag.Value().(string), tools.All())
	if err := app.DiscoverTasks(ctx); err != nil {
		return errors.Wrap(err, "discovering tasks")
	}

	taskName := parser.TaskName()
	if taskName == "" {
		return errors.Wrap(app.ListTasks(os.Stdout), "listing tasks")
	}

	task, found := lo.Find(app.Tasks, func(task application.Task) bool { return task.Name == taskName })
	if !found {
		return errors.Errorf("no task found with name %q", taskName)
	}

	if err := task.Run(ctx, parser.RemainingArgs()); err != nil {
		return errors.Wrapf(err, "running task %q", task.Name)
	}

	return nil
}
