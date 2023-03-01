package cli

import (
	"fmt"

	"github.com/samber/lo"
)

type OkArgs struct {
	Help      bool
	Version   bool
	ListTools bool
	Watches   []string
	TaskName  string
}

func (p *Parser) ParseOkArgs() (OkArgs, error) {
	args := okArgs()

	var okArgs OkArgs
	for !p.isDone() && okArgs.TaskName == "" {
		if p.currentIsFlag() {
			arg, found := lo.Find(args, func(arg okArg) bool { return arg.Match(p.current()) })
			if !found {
				return OkArgs{}, fmt.Errorf("invalid option %q", p.current())
			}

			if err := arg.Apply(p, &okArgs); err != nil {
				return OkArgs{}, err
			}
		} else {
			okArgs.TaskName = p.current()
			p.index += 1
		}
	}

	return okArgs, nil
}
