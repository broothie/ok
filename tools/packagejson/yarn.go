package packagejson

import (
	"context"
	"os/exec"

	"github.com/broothie/cob"
	"github.com/broothie/ok"
	"github.com/broothie/option"
	"github.com/samber/lo"
)

const yarnToolName = "Yarn"

func NewYarn() ok.Tool {
	return ok.Tool{
		Name:        yarnToolName,
		CommandName: "yarn",
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
						opts := option.NewOptions(
							cob.AddArgs("run", scriptName),
						)

						// Forward extra args to the script (avoid Yarn parsing them as flags).
						if len(args) > 0 {
							opts = append(opts,
								cob.AddArgs("--"),
								cob.AddArgs(args...),
							)
						}

						return opts, nil
					},
				}
			}), nil
		},
	}
}
