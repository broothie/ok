package arg

import "strings"

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

func (p *Parser) ParseNowArgs() {
	for p.argCounter < len(p.args) && p.NowArgs.TaskName == "" {
		p.processNowArg()
	}
}

func (p *Parser) processNowArg() {
	arg := p.args[p.argCounter]
	if !strings.HasPrefix(arg, "-") {
		p.NowArgs.TaskName = arg
		p.argCounter++
		return
	}

	switch arg {
	case "-h", "--help":
		p.NowArgs.Help = true
		p.argCounter++
	case "-e", "--explain":
		p.NowArgs.Explain = true
		p.argCounter++
	case "-w", "--watch":
		// TODO: check that next arg even exists
		p.NowArgs.Watches = append(p.NowArgs.Watches, p.args[p.argCounter+1])
		p.argCounter += 2
	}
}
