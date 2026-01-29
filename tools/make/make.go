package make

import (
	"context"
	"os/exec"
	"strings"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok"
	"github.com/broothie/option"
)

func New() ok.Tool {
	return ok.Tool{
		Name:        "Make",
		CommandName: "make",
		FileGlobs:   []string{"GNUmakefile", "Makefile", "makefile"},
		ProcessFile: func(ctx context.Context, filePath string, toolCfg ok.ToolConfig) ([]ok.Task, error) {
			// https://stackoverflow.com/questions/4219255/how-do-you-get-the-list-of-targets-in-a-makefile
			output, _, _, err := cob.Output(ctx, toolCfg.Executable,
				cob.AddArgs("--file", filePath),
				cob.AddArgs("--print-data-base"),      // prints the database
				cob.AddArgs("--no-builtin-variables"), // suppresses inclusion of built-in variables
				cob.AddArgs("--no-builtin-rules"),     // suppresses inclusion of built-in rules
				cob.AddArgs("--question"),             // only tests the up-to-date-status of a target (without remaking anything)
				cob.AddArgs(":"),                      // is a deliberately invalid target that is meant to ensure that no commands are executed
			)
			if err != nil && !isExitCode2(err) {
				return nil, errors.Wrapf(err, "parsing make targets from %q", filePath)
			}

			var tasks []ok.Task
			for block := range strings.SplitSeq(output.String(), "\n\n") {
				if !strings.Contains(block, "commands to execute") && !strings.Contains(block, "recipe to execute") {
					continue
				}

				lines := strings.Split(block, "\n")
				taskName := strings.TrimSuffix(lines[0], ":")
				tasks = append(tasks, ok.Task{
					Name: taskName,
					RunOptions: func(ctx context.Context, args []string, toolCfg ok.ToolConfig) (option.Options[*exec.Cmd], error) {
						return option.NewOptions(
							cob.AddArgs("--file", filePath),
							cob.AddArgs(taskName),
							cob.AddArgs(args...),
						), nil
					},
				})
			}

			return tasks, nil
		},
	}
}

func isExitCode2(err error) bool {
	exitErr := new(exec.ExitError)
	return errors.As(err, &exitErr) && exitErr.ExitCode() == 2
}
