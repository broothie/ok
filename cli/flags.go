package cli

import "fmt"

func flags() []flag {
	return []flag{
		{
			long:        "help",
			short:       'h',
			description: "Show help.",
			apply: func(parser *optionParser, options *Options) error {
				options.Help = true
				parser.advance(1)
				return nil
			},
		},
		{
			long:        "version",
			short:       'V',
			description: "Show version.",
			apply: func(parser *optionParser, options *Options) error {
				options.Version = true
				parser.advance(1)
				return nil
			},
		},
		{
			long:        "tools",
			description: "List available tools.",
			apply: func(parser *optionParser, options *Options) error {
				options.ListTools = true
				parser.advance(1)
				return nil
			},
		},
		{
			long:        "tool",
			description: "Configure a tool.",
			apply: func(parser *optionParser, options *Options) error {
				next, _ := parser.next()
				options.Tool = append(options.Tool, ParseToolOption(next))
				parser.advance(2)
				return nil
			},
		},
		{
			long:        "watch",
			short:       'w',
			description: "Glob pattern of files to watch. Can be used multiple times.",
			apply: func(parser *optionParser, options *Options) error {
				watch, present := parser.next()
				if !present {
					current, _ := parser.current()
					return fmt.Errorf("no value provided for %q", current)
				}

				options.Watches = append(options.Watches, watch.String())
				parser.advance(2)
				return nil
			},
		},
	}
}
