package cli

import (
	"errors"
	"fmt"
)

func (p *Parser) ParseOptions() (string, Options, error) {
	flagMap := make(map[string]Option)
	for _, flag := range options {
		flagMap[fmt.Sprintf("--%s", flag.Name)] = flag
		if flag.Short {
			flagMap[fmt.Sprintf("-%c", flag.Name[0])] = flag
		}
	}

	taskName := ""
	for p.argCounter < len(p.Args) && taskName == "" {
		rawArg, _ := p.current()
		if dashPrefix.MatchString(rawArg) {
			flag, flagFound := flagMap[rawArg]
			if !flagFound {
				return "", Options{}, fmt.Errorf("invalid option: '%s'", rawArg)
			}

			requiresNext := flag.ArgName != ""
			next, ok := p.peek(1)
			if requiresNext && !ok {
				return "", Options{}, fmt.Errorf("no argument provided to option '%s'", rawArg)
			}

			flag.OptionSetter(&p.options, next)
			if requiresNext {
				p.argCounter += 2
			} else {
				p.argCounter++
			}
		} else {
			taskName = rawArg
			p.argCounter++
		}
	}

	if len(p.options.Watches) > 0 && taskName == "" {
		return "", Options{}, errors.New("watches provided without task")
	}

	return taskName, p.options, nil
}
