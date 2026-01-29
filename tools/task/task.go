package task

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

type schema struct {
	Tasks []taskSchema `json:"tasks"`
}

type taskSchema struct {
	Name string `json:"name"`
	Task string `json:"task"` // Some versions also include a "task" field; accept it as fallback.
}

func New() ok.Tool {
	return ok.Tool{
		Name:        "Task",
		CommandName: "task",
		FileGlobs: []string{
			"Taskfile.yml",
			"Taskfile.yaml",
			"taskfile.yml",
			"taskfile.yaml",
		},
		ProcessFile: func(ctx context.Context, filePath string, toolCfg ok.ToolConfig) ([]ok.Task, error) {
			output, _, _, err := cob.Output(ctx, toolCfg.Executable,
				cob.AddArgs("--taskfile", filePath),
				cob.AddArgs("--list-all"),
				cob.AddArgs("--json"),
			)
			if err != nil {
				return nil, errors.Wrapf(err, "listing task tasks from %q", filePath)
			}

			var payload schema
			if err := json.NewDecoder(output).Decode(&payload); err != nil {
				return nil, errors.Wrapf(err, "parsing task task list from %q", filePath)
			}

			return lo.Map(payload.Tasks, func(tsk taskSchema, _ int) ok.Task {
				return ok.Task{
					Name: tsk.Name,
					RunOptions: func(ctx context.Context, args []string) (option.Options[*exec.Cmd], error) {
						opts := option.NewOptions(
							cob.AddArgs("--taskfile", filePath),
							cob.AddArgs(tsk.Name),
						)

						// Extra args should become CLI_ARGS, not additional tasks/flags.
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
