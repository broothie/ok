package okay

import (
	"fmt"

	"github.com/pkg/errors"
)

type Options struct {
	Help      bool
	Version   bool
	Init      string
	ListTools bool
	Watches   []string
	Stop      bool
	TaskName  string
}

type Option struct {
	Name         string
	Short        bool
	Description  string
	ArgName      string
	Stop         bool
	OptionSetter func() error
}

func (p *Parser) ParseOptions() (Options, error) {
	options := make(map[string]Option)
	for _, option := range p.availableOptions {
		options[fmt.Sprintf("--%s", option.Name)] = option
		if option.Short {
			options[fmt.Sprintf("-%c", option.Name[0])] = option
		}
	}

	for p.argCounter < len(p.rawArgs) && p.options.TaskName == "" {
		rawArg, _ := p.current()
		if dashPrefix.MatchString(rawArg) {
			option, optionFound := options[rawArg]
			if !optionFound {
				return Options{}, fmt.Errorf("invalid option: '%s'", rawArg)
			}

			if err := option.OptionSetter(); err != nil {
				return Options{}, err
			}

			if !p.options.Stop {
				p.options.Stop = option.Stop
			}
		} else {
			p.options.TaskName = rawArg
			p.argCounter++
		}
	}

	if len(p.options.Watches) > 0 && p.options.TaskName == "" {
		return Options{}, errors.New("watches provided without task")
	}

	if p.options.TaskName == "" {
		p.options.Stop = true
	}

	return p.options, nil
}

func (p *Parser) setupOptions() {
	p.availableOptions = []Option{
		{
			Name:         "help",
			Short:        true,
			Description:  "Print this help text.",
			Stop:         true,
			OptionSetter: p.helpSetter,
		},
		{
			Name:         "version",
			Short:        false,
			Description:  "Print okay version.",
			Stop:         true,
			OptionSetter: p.versionSetter,
		},
		{
			Name:         "init",
			Short:        true,
			Description:  "Initialize a tool.",
			ArgName:      "tool",
			Stop:         true,
			OptionSetter: p.initSetter,
		},
		{
			Name:         "list-tools",
			Short:        false,
			Description:  "List tools and their availability.",
			Stop:         true,
			OptionSetter: p.listToolsSetter,
		},
		{
			Name:         "watch",
			Short:        true,
			Description:  "Provide files or glob pattern to have a task run on file change.",
			ArgName:      "glob",
			Stop:         false,
			OptionSetter: p.watchSetter,
		},
	}
}

func (p *Parser) helpSetter() error {
	p.options.Help = true
	p.argCounter++
	return nil
}

func (p *Parser) versionSetter() error {
	p.options.Version = true
	p.argCounter++
	return nil
}

func (p *Parser) initSetter() error {
	toolName, ok := p.peek(1)
	if !ok {
		current, _ := p.current()
		return fmt.Errorf("no tool provided to option '%s'", current)
	}

	p.options.Init = toolName
	p.argCounter += 2
	return nil
}

func (p *Parser) listToolsSetter() error {
	p.options.ListTools = true
	p.argCounter++
	return nil
}

func (p *Parser) watchSetter() error {
	watchPattern, ok := p.peek(1)
	if !ok {
		current, _ := p.current()
		return fmt.Errorf("no watch pattern provided to option '%s'", current)
	}

	p.options.Watches = append(p.options.Watches, watchPattern)
	p.argCounter += 2
	return nil
}
