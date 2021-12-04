package cli

import (
	"time"

	"github.com/broothie/ok/config"
	"github.com/broothie/ok/logger"
)

type Options struct {
	config.Config
	Help      bool
	Version   bool
	Init      string
	ListTools bool
	Watches   []string
}

type OptionSetter func(options *Options, next string)

type Option struct {
	Name         string
	Short        bool
	Description  string
	ArgName      string
	Hidden       bool
	OptionSetter OptionSetter
}

var options = []Option{
	{
		Name:         "debug",
		Short:        false,
		Description:  "Show debug info",
		Hidden:       true,
		OptionSetter: func(options *Options, _ string) { options.Debug = true },
	},
	{
		Name:         "help",
		Short:        true,
		Description:  "Print this help text.",
		Hidden:       false,
		OptionSetter: func(options *Options, _ string) { options.Help = true },
	},
	{
		Name:         "version",
		Short:        false,
		Description:  "Print ok version.",
		Hidden:       false,
		OptionSetter: func(options *Options, _ string) { options.Version = true },
	},
	{
		Name:         "init",
		Short:        true,
		Description:  "Initialize a tool.",
		ArgName:      "tool",
		Hidden:       false,
		OptionSetter: func(options *Options, toolName string) { options.Init = toolName },
	},
	{
		Name:         "tools",
		Short:        false,
		Description:  "List tools and their availability.",
		Hidden:       false,
		OptionSetter: func(options *Options, _ string) { options.ListTools = true },
	},
	{
		Name:         "watch",
		Short:        true,
		Description:  "Provide files or glob pattern to have a task run on file change.",
		ArgName:      "glob",
		Hidden:       false,
		OptionSetter: func(options *Options, watchPattern string) { options.Watches = append(options.Watches, watchPattern) },
	},
	{
		Name:         "skip",
		Short:        false,
		Description:  "Ignore a tool.",
		ArgName:      "tool",
		Hidden:       false,
		OptionSetter: func(options *Options, toolName string) { options.SkipTools = append(options.SkipTools, toolName) },
	},
	{
		Name:        "timeout",
		Short:       false,
		Description: "Time to wait for each tool to mount.",
		ArgName:     "3s",
		Hidden:      false,
		OptionSetter: func(options *Options, timeoutString string) {
			timeout, err := time.ParseDuration(timeoutString)
			if err != nil {
				logger.Ok.Printf("failed to parse arg for --timeout: %q", timeoutString)
				return
			}

			options.Timeout = timeout
		},
	},
}
