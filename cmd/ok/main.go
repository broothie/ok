package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bobg/errors"
	"github.com/broothie/ok"
	"github.com/broothie/ok/cli"
	"github.com/broothie/ok/tools"
)

const version = "v0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func run() error {
	flags := NewFlags()

	parser := cli.NewParser(version, os.Args[1:], flags.All()...)
	if err := parser.Parse(); err != nil {
		return errors.Wrap(err, "parsing flags")
	}

	if err := os.Chdir(flags.Directory.Value().(string)); err != nil {
		return errors.Wrapf(err, "changing dir to %q", flags.Directory.Value().(string))
	}

	if flags.Help.Value().(bool) {
		return errors.Wrap(parser.WriteHelp(os.Stdout), "printing help")
	}

	ctx, cancel := context.WithTimeout(context.Background(), flags.Timeout.Value().(time.Duration))
	defer cancel()

	app := ok.New()

	if err := app.SetUpTools(ctx, tools.All()); err != nil {
		return errors.Wrap(err, "setting up tools")
	}

	if flags.Tools.Value().(bool) {
		return errors.Wrap(app.PrintTools(os.Stdout), "listing tools")
	}

	if err := app.SetUpTasks(ctx); err != nil {
		return errors.Wrap(err, "discovering tasks")
	}

	taskName := parser.TaskName()
	if taskName == "" {
		return errors.Wrap(app.ListTasks(os.Stdout), "listing tasks")
	}

	if err := app.RunTask(ctx, taskName, parser.RemainingArgs()); err != nil {
		return errors.Wrapf(err, "running task %q", taskName)
	}

	return nil
}

type Flags struct {
	Help      *cli.Flag
	Directory *cli.Flag
	Timeout   *cli.Flag
	Tools     *cli.Flag
}

func (f *Flags) All() []*cli.Flag {
	return []*cli.Flag{
		f.Help,
		f.Directory,
		f.Timeout,
		f.Tools,
	}
}

func NewFlags() *Flags {
	return &Flags{
		Help: &cli.Flag{
			Name:         "help",
			Type:         cli.FlagTypeBool,
			Help:         "Show command help.",
			Shorts:       []rune{'h'},
			DefaultValue: false,
		},
		Directory: &cli.Flag{
			Name:         "directory",
			Type:         cli.FlagTypeString,
			Help:         "Directory to run command from.",
			Shorts:       []rune{'d'},
			DefaultValue: ".",
		},
		Timeout: &cli.Flag{
			Name:         "timeout",
			Type:         cli.FlagTypeDuration,
			Help:         "Command timeout.",
			DefaultValue: time.Second,
		},
		Tools: &cli.Flag{
			Name:         "tools",
			Type:         cli.FlagTypeBool,
			Help:         "List tools.",
			DefaultValue: false,
		},
	}
}
