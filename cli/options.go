package cli

import (
	"fmt"

	"github.com/samber/lo"
)

type Options struct {
	Tool      []ToolOptions
	Help      bool
	Version   bool
	ListTools bool
	InitTool  string
	Watches   []string
	TaskName  string
}

// Options can only be called once.
func (cli *CLI) Options() (Options, error) {
	var options Options
	for cli.parser.hasArgsLeft() {
		current, _ := cli.parser.current()
		if current.isFlag() {
			flag, present := lo.Find(cli.flags, func(flag flag) bool { return flag.isMatch(current) })
			if !present {
				return Options{}, fmt.Errorf("unknown flag %q", current)
			}

			if err := flag.apply(cli.parser, &options); err != nil {
				return Options{}, err
			}
		} else {
			options.TaskName = current.String()
			cli.parser.advance(1)
			break
		}
	}

	return options, nil
}
