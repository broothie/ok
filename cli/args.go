package cli

import "fmt"

func okArgs() []okArg {
	return []okArg{
		{
			Name:        "help",
			Short:       'h',
			Description: "Show help.",
			Apply: func(parser *Parser, okArgs *OkArgs) error {
				okArgs.Help = true
				parser.index += 1
				return nil
			},
		},
		{
			Name:        "version",
			Short:       'V',
			Description: "Show version.",
			Apply: func(parser *Parser, okArgs *OkArgs) error {
				okArgs.Version = true
				parser.index += 1
				return nil
			},
		},
		{
			Name:        "tools",
			Description: "List available tools.",
			Apply: func(parser *Parser, okArgs *OkArgs) error {
				okArgs.ListTools = true
				parser.index += 1
				return nil
			},
		},
		{
			Name:        "watch",
			Short:       'w',
			Description: "Glob pattern of files to watch. Can be used multiple times.",
			Apply: func(parser *Parser, okArgs *OkArgs) error {
				watch, present := parser.peek(1)
				if !present {
					return fmt.Errorf("no value provided for %q", parser.current())
				}

				okArgs.Watches = append(okArgs.Watches, watch)
				parser.index += 2
				return nil
			},
		},
	}
}
