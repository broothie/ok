package mise

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok"
	"github.com/broothie/option"
	"github.com/samber/lo"
)

type taskSchema struct {
	Name string `json:"name"`
}

func New() ok.Tool {
	return ok.Tool{
		Name:        "Mise",
		CommandName: "mise",
		FileGlobs:   []string{"mise.toml", ".mise.toml", ".mise/config.toml"},
		ProcessFile: func(ctx context.Context, filePath string, toolCfg ok.ToolConfig) ([]ok.Task, error) {
			output, _, _, err := cob.Output(ctx, toolCfg.Executable,
				cob.AddArgs("tasks", "--json"),
			)
			if err != nil {
				return nil, errors.Wrapf(err, "listing mise tasks from %q", filePath)
			}

			var payload []taskSchema
			if err := json.NewDecoder(output).Decode(&payload); err != nil {
				return nil, errors.Wrapf(err, "parsing mise tasks from %q", filePath)
			}

			return lo.Map(payload, func(task taskSchema, _ int) ok.Task {
				return ok.Task{
					Name: task.Name,
					RunOptions: func(ctx context.Context, args []string, toolCfg ok.ToolConfig) (option.Options[*exec.Cmd], error) {
						opts := option.NewOptions(
							cob.AddArgs("run"),
							cob.AddArgs(task.Name),
						)

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
