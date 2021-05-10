package arg

import (
	"fmt"
	"strings"
)

type Options struct {
	Help      bool
	Init      string
	ListTools bool
	Watches   []string
	TaskName  string
}

func (o Options) WillRunTask() bool {
	if o.Help || o.ListTools || o.Init != "" {
		return false
	}

	return o.TaskName != ""
}

func (p *Parser) ParseOptions() error {
	for p.argCounter < len(p.args) && p.Options.TaskName == "" {
		if err := p.processOption(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) processOption() error {
	arg, _ := p.current()
	if !strings.HasPrefix(arg, "-") {
		p.Options.TaskName = arg
		p.argCounter++
		return nil
	}

	switch arg {
	case "-h", "--help":
		p.Options.Help = true
		p.argCounter++

	case "--list-tools":
		p.Options.ListTools = true
		p.argCounter++

	case "-i", "--init":
		taskName, ok := p.peek(1)
		if !ok {
			return fmt.Errorf("no task provided to '%s'", arg)
		}

		p.Options.Init = taskName
		p.argCounter += 1

	case "-w", "--watch":
		watchPattern, ok := p.peek(1)
		if !ok {
			return fmt.Errorf("no pattern provided to '%s'", arg)
		}

		p.Options.Watches = append(p.Options.Watches, watchPattern)
		p.argCounter += 2

	default:
		return fmt.Errorf("invalid arg '%s'", arg)
	}

	return nil
}
