package cli

import (
	"fmt"

	"github.com/broothie/ok/argument"
	"github.com/broothie/ok/parameter"
	"github.com/pkg/errors"
)

func (p *Parser) ParseTaskArgs(parameters parameter.Parameters) (argument.Arguments, error) {
	var args argument.Arguments

	requiredIndex := 0
	for !p.isDone() {
		if !p.currentIsFlag() {
			param, present := parameters.Required().Get(requiredIndex)
			if !present {
				return nil, errors.New("too many arguments")
			}

			args = append(args, argument.Argument{Parameter: param, Value: p.current()})
			p.index += 1
			requiredIndex += 1
		} else {
			if p.currentIsShortFlag() {
				return nil, fmt.Errorf("short flags unsupported for task args: %q", p.current())
			}

			name := p.currentDashless()
			param, present := parameters.Optional().Find(func(p parameter.Parameter) bool { return p.Name == name })
			if !present {
				return nil, fmt.Errorf("no parameter with name %q", name)
			}

			value, present := p.peek(1)
			if !present {
				return nil, fmt.Errorf("no value provided for %q", p.current())
			}

			args = append(args, argument.Argument{Parameter: param, Value: value})
			p.index += 2
		}
	}

	if len(args.Required()) < len(parameters.Required()) {
		return nil, fmt.Errorf("missing required args (given %d, expected %d)", len(args.Required()), len(parameters.Required()))
	}

	return args, nil
}
