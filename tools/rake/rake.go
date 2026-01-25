package rake

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
		Name:        "Rake",
		CommandName: "rake",
		FileGlobs:   []string{"Rakefile", "rakefile", "Rakefile.rb", "rakefile.rb"},
		ProcessFile: func(ctx context.Context, filePath string) ([]ok.Task, error) {
			output, _, _, err := cob.Output(ctx, "rake",
				cob.AddArgs("--rakefile", filePath),
				cob.AddArgs("--tasks"),
				cob.AddArgs("--all"),
			)
			if err != nil {
				return nil, errors.Wrapf(err, "listing rake tasks from %q", filePath)
			}

			var tasks []ok.Task
			for line := range strings.SplitSeq(output.String(), "\n") {
				line = strings.TrimSpace(line)
				if !strings.HasPrefix(line, "rake ") {
					continue
				}

				// Format: "rake task_name[args]  # Description"
				line = strings.TrimPrefix(line, "rake ")
				taskName, _, _ := strings.Cut(line, " ")
				taskName, _, _ = strings.Cut(taskName, "[") // strip arg placeholders

				tasks = append(tasks, ok.Task{
					Name: taskName,
					RunOptions: func(ctx context.Context, args []string) (option.Options[*exec.Cmd], error) {
						return option.NewOptions(
							cob.AddArgs("--rakefile", filePath),
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
