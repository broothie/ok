package task

import (
	"fmt"

	"github.com/broothie/ok/arg"
	"github.com/samber/lo"
)

func (p Parameters) Parse(parser *arg.Parser) (Arguments, error) {
	positionalParams := lo.Filter(p, func(param Parameter, _ int) bool { return param.IsRequired() })
	keywordParams := lo.Filter(p, func(param Parameter, _ int) bool { return param.IsOptional() })

	positionalIndex := 0
	var args Arguments
	for parser.HasArgsLeft() {
		current, _ := parser.Current()
		if current.IsFlag() {
			param, found := lo.Find(keywordParams, func(param Parameter) bool { return param.Name == current.Dashless() })
			if !found {
				return nil, fmt.Errorf("unknown task arg %q", current)
			}

			value, found := parser.Next()
			if !found {
				return nil, fmt.Errorf("no value provided for %q", current)
			}

			args = append(args, Argument{
				Parameter: param,
				Value:     value.String(),
			})

			parser.Advance(2)
		} else {
			if positionalIndex >= len(positionalParams) {
				return nil, fmt.Errorf("")
			}

			param := positionalParams[positionalIndex]
			args = append(args, Argument{
				Parameter: param,
				Value:     current.String(),
			})

			positionalIndex += 1
			parser.Advance(1)
		}
	}

	if positionalIndex < len(positionalParams) {
		return nil, fmt.Errorf("missing required args (given %d, expected %d)", positionalIndex, len(positionalParams))
	}

	return args, nil
}
