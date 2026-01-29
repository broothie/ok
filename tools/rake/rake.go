package rake

import (
	"context"
	"fmt"
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
		ProcessFile: func(ctx context.Context, filePath string, toolCfg ok.ToolConfig) ([]ok.Task, error) {
			output, _, _, err := cob.Output(ctx, toolCfg.Executable,
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
					RunOptions: func(ctx context.Context, args []string, toolCfg ok.ToolConfig) (option.Options[*exec.Cmd], error) {
						taskArgs, rakeVars, err := splitRakeArgs(args)
						if err != nil {
							return nil, errors.Wrap(err, "parsing rake args")
						}

						invocation := taskName
						if len(taskArgs) > 0 {
							invocation = fmt.Sprintf("%s[%s]", invocation, strings.Join(taskArgs, ","))
						}

						return option.NewOptions(
							cob.AddArgs("--rakefile", filePath),
							cob.AddArgs(invocation),
							cob.AddArgs(rakeVars...),
						), nil
					},
				})
			}

			return tasks, nil
		},
	}
}

func splitRakeArgs(args []string) (taskArgs []string, rakeVars []string, err error) {
	for _, arg := range args {
		// Common rake CLI usage: `rake task FOO=bar`
		// Treat KEY=VALUE tokens as rake vars, not task args.
		if !strings.HasPrefix(arg, "-") && strings.Contains(arg, "=") {
			rakeVars = append(rakeVars, arg)
			continue
		}

		// Rake task args are comma-separated inside `task[ ... ]`.
		// We can't represent literal commas in a single argument unambiguously.
		if strings.ContainsAny(arg, "],") {
			return nil, nil, errors.Errorf("rake task argument %q contains ',' or ']' which cannot be encoded in task[arg1,arg2] form", arg)
		}

		taskArgs = append(taskArgs, arg)
	}

	return taskArgs, rakeVars, nil
}
