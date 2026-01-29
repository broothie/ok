package packagejson

import (
	"context"
	"os/exec"

	"github.com/broothie/cob"
	"github.com/broothie/ok"
	"github.com/broothie/option"
	"github.com/samber/lo"
)

func NewNPM() ok.Tool {
	return ok.Tool{
		Name:        "NPM",
		CommandName: "npm",
		FileGlobs:   []string{packageJSONFileName},
		ProcessFile: func(ctx context.Context, filePath string, toolCfg ok.ToolConfig) ([]ok.Task, error) {
			schema, err := read(filePath)
			if err != nil {
				return nil, err
			}

			return lo.Map(lo.Keys(schema.Scripts), func(scriptName string, _ int) ok.Task {
				return ok.Task{
					Name: scriptName,
					RunOptions: func(ctx context.Context, args []string) (option.Options[*exec.Cmd], error) {
						return option.NewOptions(
							cob.AddArgs("run", scriptName),
							cob.AddArgs(args...),
						), nil
					},
				}
			}), nil
		},
	}
}
