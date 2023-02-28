package cli

import "fmt"

type OkArgs struct {
	Help      bool
	ListTools bool
	TaskName  string
}

func (p *Parser) ParseOkArgs() (OkArgs, error) {
	var okArgs OkArgs
	for !p.isDone() && okArgs.TaskName == "" {
		switch p.current() {
		case "-h", "--help":
			okArgs.Help = true
			p.index += 1

		case "--tools":
			okArgs.ListTools = true
			p.index += 1

		default:
			if p.currentIsFlag() {
				return OkArgs{}, fmt.Errorf("invalid flag %q", p.current())
			}

			okArgs.TaskName = p.current()
			p.index += 1
		}
	}

	return okArgs, nil
}
