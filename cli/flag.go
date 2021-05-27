package cli

type OptionSetter func(options *Options, next string)

type Flag struct {
	Name         string
	Short        bool
	Description  string
	ArgName      string
	Hidden       bool
	Halt         bool
	OptionSetter OptionSetter
}

var Flags = []Flag{
	{
		Name:         "debug",
		Short:        false,
		Description:  "Show debug info",
		Hidden:       true,
		Halt:         false,
		OptionSetter: func(options *Options, _ string) { options.Debug = true },
	},
	{
		Name:         "help",
		Short:        true,
		Description:  "Print this help text.",
		Halt:         true,
		OptionSetter: func(options *Options, _ string) { options.Help = true },
	},
	{
		Name:         "version",
		Short:        false,
		Description:  "Print ok version.",
		Halt:         true,
		OptionSetter: func(options *Options, _ string) { options.Version = true },
	},
	{
		Name:         "init",
		Short:        true,
		Description:  "Initialize a tool.",
		ArgName:      "tool",
		Halt:         true,
		OptionSetter: func(options *Options, toolName string) { options.Init = toolName },
	},
	{
		Name:         "tools",
		Short:        false,
		Description:  "List tools and their availability.",
		Halt:         true,
		OptionSetter: func(options *Options, _ string) { options.ListTools = true },
	},
	{
		Name:         "watch",
		Short:        true,
		Description:  "Provide files or glob pattern to have a task run on file change.",
		ArgName:      "glob",
		Halt:         false,
		OptionSetter: func(options *Options, watchPattern string) { options.Watches = append(options.Watches, watchPattern) },
	},
	{
		Name:         "skip",
		Description:  "Ignore a tool.",
		ArgName:      "tool",
		Halt:         false,
		OptionSetter: func(options *Options, toolName string) { options.SkipTools = append(options.SkipTools, toolName) },
	},
	{
		Name:         "zsh-autocomplete",
		Description:  "Outputs zsh autocomplete script.",
		Hidden:       true,
		Halt:         true,
		OptionSetter: func(options *Options, toolName string) { options.ZshAutocomplete = true },
	},
}
