package make

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok/tool"
)

func NewMake() tool.Tool {
	return tool.Tool{
		Name:        "Make",
		CommandName: "make",
		FileGlobs:   []string{"GNUmakefile", "Makefile", "makefile"},
		ParseFile: func(ctx context.Context, filePath string) ([]tool.Task, error) {
			// https://stackoverflow.com/questions/4219255/how-do-you-get-the-list-of-targets-in-a-makefile
			output, _, _, err := cob.Output(ctx, "make",
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

			var tasks []tool.Task
			for block := range strings.SplitSeq(output.String(), "\n\n") {
				if !strings.Contains(block, "commands to execute") {
					continue
				}

				lines := strings.Split(block, "\n")
				taskName := strings.TrimSuffix(lines[0], ":")
				tasks = append(tasks, tool.Task{
					Name: taskName,
					Run: func(ctx context.Context, args []string) error {
						_, err := cob.Run(ctx, "make",
							cob.AddArgs("--file", filePath),
							cob.AddArgs(taskName),
							cob.AddArgs(args...),
							cob.SetStdin(os.Stdin),
							cob.SetStdout(os.Stdout),
							cob.SetStderr(os.Stderr),
						)

						return errors.Wrapf(err, "running make target %q", taskName)
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
