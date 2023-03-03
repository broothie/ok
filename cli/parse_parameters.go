package cli

import (
	"fmt"

	"github.com/broothie/ok/task"
	"github.com/samber/lo"
)

func (cli *CLI) ParseParameters(params task.Parameters) (task.Arguments, error) {
	positionalParams := lo.Filter(params, func(param task.Parameter, _ int) bool { return param.IsPositional() })
	keywordParams := lo.Filter(params, func(param task.Parameter, _ int) bool { return param.IsKeyword() })

	positionalIndex := 0
	var args task.Arguments
	for cli.parser.hasArgsLeft() {
		current, _ := cli.parser.current()

		// If splat, eat all args
		if params.IsSplat() {
			args = append(args, task.Argument{Parameter: params[0], Value: current.String()})
			cli.parser.advance(1)
			continue
		}

		if current.isFlag() {
			param, found := lo.Find(keywordParams, func(param task.Parameter) bool { return param.Name == current.dashless() })
			if !found {
				return nil, fmt.Errorf("unknown task arg %q", current)
			}

			value, found := cli.parser.next()
			if !found {
				return nil, fmt.Errorf("no value provided for %q", current)
			}

			args = append(args, task.Argument{
				Parameter: param,
				Value:     value.String(),
			})

			cli.parser.advance(2)
		} else {
			if positionalIndex >= len(positionalParams) {
				return nil, fmt.Errorf("too many positional args (given %d, expected %d)", positionalIndex+1, len(positionalParams))
			}

			param := positionalParams[positionalIndex]
			args = append(args, task.Argument{
				Parameter: param,
				Value:     current.String(),
			})

			positionalIndex += 1
			cli.parser.advance(1)
		}
	}

	if positionalIndex < len(positionalParams) && !params.IsSplat() {
		return nil, fmt.Errorf("missing required args (given %d, expected %d)", positionalIndex, len(positionalParams))
	}

	return args, nil
}
