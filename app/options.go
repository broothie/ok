package app

import (
	"fmt"

	"github.com/broothie/ok/arg"
	"github.com/samber/lo"
)

type Options struct {
	Help      bool
	Version   bool
	ListTools bool
	Watches   []string
	TaskName  string
}

func (app *App) Options(parser *arg.Parser) (Options, error) {
	flags := app.Flags()

	var options Options
	for !parser.IsExhausted() {
		current, _ := parser.Current()
		if current.IsFlag() {
			flag, present := lo.Find(flags, func(flag Flag) bool { return flag.IsMatch(current) })
			if !present {
				return Options{}, fmt.Errorf("unrecognized flag %q", current)
			}

			if err := flag.Apply(parser, &options); err != nil {
				return Options{}, err
			}
		} else {
			options.TaskName = current.String()
			parser.Advance(1)
			break
		}
	}

	return options, nil
}
