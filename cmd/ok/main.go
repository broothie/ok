package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/bobg/errors"
	"github.com/broothie/ok"
	"github.com/broothie/ok/cli"
	"github.com/broothie/ok/tools"
	"github.com/joho/godotenv"
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
	flags := NewFlags()

	parser := cli.NewParser(version, os.Args[1:], flags.All()...)
	if err := parser.Parse(); err != nil {
		return errors.Wrap(err, "parsing flags")
	}

	if flags.Version.Value().(bool) {
		fmt.Printf("ok %s\n", version)
		return nil
	}

	if flags.Debug.Value().(bool) {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if flags.LoadDotEnv.Value().(bool) {
		if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
			return errors.Wrap(err, "loading .env")
		}
	}

	if err := os.Chdir(flags.Directory.Value().(string)); err != nil {
		return errors.Wrapf(err, "changing dir to %q", flags.Directory.Value().(string))
	}

	if flags.Help.Value().(bool) {
		return errors.Wrap(parser.WriteHelp(os.Stdout), "printing help")
	}

	ctx, cancel := context.WithTimeout(context.Background(), flags.Timeout.Value().(time.Duration))
	defer cancel()

	tls := tools.All()
	if filterTools := flags.FilterTools.Value().(string); filterTools != "" {
		selectTools := strings.Split(filterTools, ",")
		tls = lo.Filter(tools.All(), func(tl ok.Tool, _ int) bool {
			return lo.ContainsBy(selectTools, func(toolName string) bool { return strings.EqualFold(toolName, tl.Name) })
		})
	}

	app := ok.New()
	if err := app.SetUpTools(ctx, tls); err != nil {
		return errors.Wrap(err, "setting up tools")
	}

	if flags.ListTools.Value().(bool) {
		return errors.Wrap(app.ListTools(os.Stdout), "listing tools")
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
	Version     *cli.Flag
	Help        *cli.Flag
	Directory   *cli.Flag
	Timeout     *cli.Flag
	FilterTools *cli.Flag
	ListTools   *cli.Flag
	LoadDotEnv  *cli.Flag
	Debug       *cli.Flag
}

func (f *Flags) All() []*cli.Flag {
	return []*cli.Flag{
		f.Version,
		f.Help,
		f.Directory,
		f.Timeout,
		f.FilterTools,
		f.ListTools,
		f.LoadDotEnv,
		f.Debug,
	}
}

func NewFlags() *Flags {
	return &Flags{
		Version: &cli.Flag{
			Name:         "version",
			Type:         cli.FlagTypeBool,
			Help:         "Print command version.",
			Shorts:       []rune{'V'},
			DefaultValue: false,
		},
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
		FilterTools: &cli.Flag{
			Name:         "filter-tools",
			Type:         cli.FlagTypeString,
			Help:         "Filter tools by case-insensitive name. Use commas for multiple values",
			Aliases:      []string{"ft"},
			DefaultValue: "",
		},
		ListTools: &cli.Flag{
			Name:         "list-tools",
			Type:         cli.FlagTypeBool,
			Help:         "List tools.",
			DefaultValue: false,
		},
		LoadDotEnv: &cli.Flag{
			Name:         "load-dot-env",
			Type:         cli.FlagTypeBool,
			Help:         "Pick up local .env files.",
			DefaultValue: true,
		},
		Debug: &cli.Flag{
			Name:         "debug",
			Type:         cli.FlagTypeBool,
			Help:         "Output debug logs.",
			DefaultValue: false,
		},
	}
}
