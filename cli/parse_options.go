package cli

import (
	"errors"
	"fmt"
)

func (p *Parser) ParseFlags() (Options, error) {
	flagMap := make(map[string]Flag)
	for _, flag := range Flags {
		flagMap[fmt.Sprintf("--%s", flag.Name)] = flag
		if flag.Short {
			flagMap[fmt.Sprintf("-%c", flag.Name[0])] = flag
		}
	}

	for p.argCounter < len(p.Args) && p.config.TaskName == "" {
		rawArg, _ := p.current()
		if dashPrefix.MatchString(rawArg) {
			flag, flagFound := flagMap[rawArg]
			if !flagFound {
				return Options{}, fmt.Errorf("invalid option: '%s'", rawArg)
			}

			requiresNext := flag.ArgName != ""
			next, ok := p.peek(1)
			if requiresNext && !ok {
				return Options{}, fmt.Errorf("no argument provided to option '%s'", rawArg)
			}

			flag.OptionSetter(&p.config, next)
			if !p.config.Halt {
				p.config.Halt = flag.Halt
			}

			if requiresNext {
				p.argCounter += 2
			} else {
				p.argCounter++
			}
		} else {
			p.config.TaskName = rawArg
			p.argCounter++
		}
	}

	if len(p.config.Watches) > 0 && p.config.TaskName == "" {
		return Options{}, errors.New("watches provided without task")
	}

	if p.config.TaskName == "" {
		p.config.Halt = true
	}

	return p.config, nil
}
