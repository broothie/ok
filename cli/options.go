package cli

import "fmt"

type OptionSetter func(parser *Parser) error

type Option struct {
	Name         string
	Short        bool
	Description  string
	ArgName      string
	Stop         bool
	OptionSetter OptionSetter
}

var Options = []Option{
	{
		Name:        "help",
		Short:       true,
		Description: "Print this help text.",
		Stop:        true,
		OptionSetter: func(parser *Parser) error {
			parser.options.Help = true
			parser.argCounter++
			return nil
		},
	},
	{
		Name:        "version",
		Short:       false,
		Description: "Print ok version.",
		Stop:        true,
		OptionSetter: func(parser *Parser) error {
			parser.options.Version = true
			parser.argCounter++
			return nil
		},
	},
	{
		Name:        "init",
		Short:       true,
		Description: "Initialize a tool.",
		ArgName:     "tool",
		Stop:        true,
		OptionSetter: func(parser *Parser) error {
			toolName, ok := parser.peek(1)
			if !ok {
				current, _ := parser.current()
				return fmt.Errorf("no tool provided to option '%s'", current)
			}

			parser.options.Init = toolName
			parser.argCounter += 2
			return nil
		},
	},
	{
		Name:        "tools",
		Short:       false,
		Description: "List tools and their availability.",
		Stop:        true,
		OptionSetter: func(parser *Parser) error {
			parser.options.ListTools = true
			parser.argCounter++
			return nil
		},
	},
	{
		Name:        "watch",
		Short:       true,
		Description: "Provide files or glob pattern to have a task run on file change.",
		ArgName:     "glob",
		Stop:        false,
		OptionSetter: func(parser *Parser) error {
			watchPattern, ok := parser.peek(1)
			if !ok {
				current, _ := parser.current()
				return fmt.Errorf("no watch pattern provided to option '%s'", current)
			}

			parser.options.Watches = append(parser.options.Watches, watchPattern)
			parser.argCounter += 2
			return nil
		},
	},
}
