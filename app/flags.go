package app

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/broothie/ok"
	"github.com/broothie/ok/arg"
	"github.com/pkg/errors"
)

type applyFunc func(*arg.Parser, *Options) error

type Flag struct {
	Long        string
	Short       rune
	Description string
	Apply       applyFunc
}

func (f Flag) IsMatch(arg arg.Token) bool {
	if !arg.IsFlag() {
		return false
	} else if arg.IsLongFlag() {
		return f.Long == arg.Dashless()
	} else {
		return string(f.Short) == arg.Dashless()
	}
}

func (f Flag) HasShort() bool {
	return f.Short != 0
}

func (app *App) PrintHelp() error {
	fmt.Printf("ok %s\n", ok.Version())
	fmt.Printf("\n")
	fmt.Printf("Usage:\n")
	fmt.Printf("  ok [OPTIONS] <TASK> [TASK ARGS]\n")
	fmt.Printf("\n")
	fmt.Printf("Options:\n")

	var rows []string
	for _, flag := range app.Flags() {
		short := ""
		if flag.HasShort() {
			short = fmt.Sprintf("-%c", flag.Short)
		}

		row := []string{fmt.Sprintf("  %s", short), fmt.Sprintf("--%s", flag.Long), flag.Description}
		rows = append(rows, strings.Join(row, "\t"))
	}

	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, strings.Join(rows, "\n")); err != nil {
		return errors.Wrap(err, "failed to write rows to table")
	}

	if err := table.Flush(); err != nil {
		return errors.Wrap(err, "failed to write table")
	}

	return nil
}

func (app *App) Flags() []Flag {
	return []Flag{
		{
			Long:        "help",
			Short:       'h',
			Description: "Show help.",
			Apply: func(parser *arg.Parser, options *Options) error {
				options.Help = true
				parser.Advance(1)
				return nil
			},
		},
		{
			Long:        "version",
			Short:       'V',
			Description: "Show version.",
			Apply: func(parser *arg.Parser, options *Options) error {
				options.Version = true
				parser.Advance(1)
				return nil
			},
		},
		{
			Long:        "tools",
			Description: "List available tools.",
			Apply: func(parser *arg.Parser, options *Options) error {
				options.ListTools = true
				parser.Advance(1)
				return nil
			},
		},
		{
			Long:        "watch",
			Short:       'w',
			Description: "Glob pattern of files to watch. Can be used multiple times.",
			Apply: func(parser *arg.Parser, options *Options) error {
				watch, present := parser.Next()
				if !present {
					current, _ := parser.Current()
					return fmt.Errorf("no value provided for %q", current)
				}

				options.Watches = append(options.Watches, watch.String())
				parser.Advance(2)
				return nil
			},
		},
	}
}
