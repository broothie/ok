package arg

import (
	"errors"
	"fmt"
	"strings"
)

type Now struct {
	Help     bool
	Explain  bool
	Watches  []string
	TaskName string
}

func (n Now) WillRunTask() bool {
	if n.Help || n.Explain {
		return false
	}

	return n.TaskName != ""
}

func (p *Parser) ParseNowArgs() error {
	for p.argCounter < len(p.args) && p.NowArgs.TaskName == "" {
		if err := p.processNowArg(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) processNowArg() error {
	arg, _ := p.current()
	if !strings.HasPrefix(arg, "-") {
		p.NowArgs.TaskName = arg
		p.argCounter++
		return nil
	}

	switch arg {
	case "-h", "--help":
		p.NowArgs.Help = true
		p.argCounter++

	case "-e", "--explain":
		p.NowArgs.Explain = true
		p.argCounter++

	case "-w", "--watch":
		watchPattern, ok := p.peek(1)
		if !ok {
			return errors.New("no watch pattern provided to watch")
		}

		p.NowArgs.Watches = append(p.NowArgs.Watches, watchPattern)
		p.argCounter += 2

	default:
		return fmt.Errorf("invalid arg '%s'", arg)
	}

	return nil
}
