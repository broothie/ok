package okay

import (
	"fmt"
	"strings"

	"github.com/broothie/okay/task"
)

func (p *Parser) ParseArgs(params task.Parameters) (task.Args, error) {
	if params.Forward {
		return task.Args{Forwards: p.rawArgs[p.argCounter:]}, nil
	}

	args := task.Args{Keyword: make(map[string]task.Arg)}

	for p.argCounter < len(p.rawArgs) {
		rawArg, _ := p.current()
		if dashPrefix.MatchString(rawArg) {
			// Keyword
			argSansDash := dashPrefix.ReplaceAllString(rawArg, "")
			param, paramFound := params.KeywordAt(argSansDash)
			if !paramFound {
				return task.Args{}, fmt.Errorf("invalid keyword arg '%s'", rawArg)
			}

			if param.Type == task.Bool {
				value, defaultExists := param.Default.(bool)
				if defaultExists {
					value = !value
				} else {
					value = true
				}

				args.Keyword[param.Name] = task.Arg{Parameter: param, Value: value}
				p.argCounter++
				continue
			}

			valueArg, valuePresent := p.peek(1)
			if !valuePresent {
				return task.Args{}, fmt.Errorf("no value provided to keyword arg '%s'", rawArg)
			}

			arg, err := processArgWithParam(valueArg, param)
			if err != nil {
				return task.Args{}, err
			}

			args.Keyword[param.Name] = arg
			p.argCounter += 2
		} else {
			// Positional
			param, paramPresent := params.PositionalAt(len(args.Positional))
			if !paramPresent {
				return task.Args{}, fmt.Errorf("too many positional args provided, expected max of %d", len(params.PositionalRequired)+len(params.PositionalOptional))
			}

			arg, err := processArgWithParam(rawArg, param)
			if err != nil {
				return task.Args{}, err
			}

			args.Positional = append(args.Positional, arg)
			p.argCounter++
		}
	}

	if len(args.Positional) < len(params.PositionalRequired) {
		return task.Args{}, missingPositionalError(params.PositionalRequired, args.Positional)
	}

	if len(args.Keyword) < len(params.KeywordRequired) {
		return task.Args{}, missingKeywordArgError(params.KeywordRequired, args.Keyword)
	}

	return args, nil
}

//func (p *Parser) forward(params task.Parameters) (task.Args, error) {
//	for p.argCounter < len(p.rawArgs) {
//
//	}
//}

func processArgWithParam(rawArg string, param task.Parameter) (task.Arg, error) {
	processed, err := param.Type.Parse(rawArg)
	if err != nil {
		return task.Arg{}, err
	}

	return task.Arg{Parameter: param, Value: processed}, nil
}

func missingPositionalError(params []task.Parameter, args []task.Arg) error {
	missingArgs := make([]string, 0, len(params))
	for _, param := range params[len(args):] {
		missingArgs = append(missingArgs, param.Name)
	}

	return fmt.Errorf("missing positional args: [%s]", strings.Join(missingArgs, ", "))
}

func missingKeywordArgError(params []task.Parameter, args map[string]task.Arg) error {
	missingArgs := make([]string, 0, len(args))
	for _, param := range params {
		if _, argPresent := args[param.Name]; !argPresent {
			missingArgs = append(missingArgs, param.Name)
		}
	}

	return fmt.Errorf("missing keyword args: [%s]", strings.Join(missingArgs, ", "))
}
