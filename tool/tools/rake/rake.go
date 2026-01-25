package rake

import (
	"context"
	"os"
	"strings"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok/tool"
)

func NewRake() tool.Tool {
	return tool.Tool{
		Name:        "Rake",
		CommandName: "rake",
		FileGlobs:   []string{"Rakefile", "rakefile", "Rakefile.rb", "rakefile.rb"},
		ParseFile: func(ctx context.Context, filePath string) ([]tool.Task, error) {
			output, _, _, err := cob.Output(ctx, "rake",
				cob.AddArgs("--rakefile", filePath),
				cob.AddArgs("--tasks"),
				cob.AddArgs("--all"),
			)
			if err != nil {
				return nil, errors.Wrapf(err, "listing rake tasks from %q", filePath)
			}

			var tasks []tool.Task
			for line := range strings.SplitSeq(output.String(), "\n") {
				line = strings.TrimSpace(line)
				if !strings.HasPrefix(line, "rake ") {
					continue
				}

				// Format: "rake task_name[args]  # Description"
				line = strings.TrimPrefix(line, "rake ")
				taskName, _, _ := strings.Cut(line, " ")
				taskName, _, _ = strings.Cut(taskName, "[") // strip arg placeholders

				tasks = append(tasks, tool.Task{
					Name: taskName,
					Run: func(ctx context.Context, args []string) error {
						_, err := cob.Run(ctx, "rake",
							cob.AddArgs("--rakefile", filePath),
							cob.AddArgs(taskName),
							cob.AddArgs(args...),
							cob.SetStdin(os.Stdin),
							cob.SetStdout(os.Stdout),
							cob.SetStderr(os.Stderr),
						)

						return errors.Wrapf(err, "running rake task %q", taskName)
					},
				})
			}

			return tasks, nil
		},
	}
}
