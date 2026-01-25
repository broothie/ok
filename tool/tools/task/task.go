package task

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok/tool"
	"github.com/broothie/option"
)

func New() tool.Tool {
	return tool.Tool{
		Name:        "Task",
		CommandName: "task",
		FileGlobs: []string{
			"Taskfile.yml",
			"Taskfile.yaml",
			"taskfile.yml",
			"taskfile.yaml",
		},
		ParseFile: func(ctx context.Context, filePath string) ([]tool.Task, error) {
			output, _, _, err := cob.Output(ctx, "task",
				cob.AddArgs("--taskfile", filePath),
				cob.AddArgs("--list-all"),
				cob.AddArgs("--json"),
			)
			if err != nil {
				return nil, errors.Wrapf(err, "listing task tasks from %q", filePath)
			}

			var payload struct {
				Tasks []struct {
					Name string `json:"name"`
					Task string `json:"task"` // Some versions also include a "task" field; accept it as fallback.
				} `json:"tasks"`
			}
			if err := json.NewDecoder(output).Decode(&payload); err != nil {
				return nil, errors.Wrapf(err, "parsing task task list from %q", filePath)
			}

			var tasks []tool.Task
			for _, task := range payload.Tasks {
				taskName := task.Name
				if taskName == "" {
					taskName = task.Task
				}

				tasks = append(tasks, tool.Task{
					Name: taskName,
					RunOptions: func(ctx context.Context, args []string) (option.Options[*exec.Cmd], error) {
						opts := option.NewOptions(
							cob.AddArgs("--taskfile", filePath),
							cob.AddArgs(taskName),
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
				})
			}

			return tasks, nil
		},
	}
}
