package cli

import (
	"errors"
	"fmt"

	"github.com/broothie/ok/ok"
)

func (p *Parser) ParseOptions() (ok.Options, error) {
	optionMap := make(map[string]Option)
	for _, option := range Options {
		optionMap[fmt.Sprintf("--%s", option.Name)] = option
		if option.Short {
			optionMap[fmt.Sprintf("-%c", option.Name[0])] = option
		}
	}

	for p.argCounter < len(p.Args) && p.options.TaskName == "" {
		rawArg, _ := p.current()
		if dashPrefix.MatchString(rawArg) {
			option, optionFound := optionMap[rawArg]
			if !optionFound {
				return ok.Options{}, fmt.Errorf("invalid option: '%s'", rawArg)
			}

			if err := option.OptionSetter(p); err != nil {
				return ok.Options{}, err
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
		return ok.Options{}, errors.New("watches provided without task")
	}

	if p.options.TaskName == "" {
		p.options.Stop = true
	}

	return p.options, nil
}
