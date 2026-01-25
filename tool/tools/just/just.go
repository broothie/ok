package just

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok"
	"github.com/broothie/option"
)

func New() ok.Tool {
	return ok.Tool{
		Name:        "Just",
		CommandName: "just",
		FileGlobs:   []string{"Justfile", "justfile"},
		ProcessFile: func(ctx context.Context, filePath string) ([]ok.Task, error) {
			output, _, _, err := cob.Output(ctx, "just",
				cob.AddArgs("--justfile", filePath),
				cob.AddArgs("--dump", "--dump-format", "json"),
			)
			if err != nil {
				return nil, errors.Wrapf(err, "listing just recipes from %q", filePath)
			}

			var justfile struct {
				Recipes map[string]struct {
					Name string `json:"name"`
				} `json:"recipes"`
			}
			if err := json.NewDecoder(output).Decode(&justfile); err != nil {
				return nil, errors.Wrapf(err, "parsing just recipes from %q", filePath)
			}

			var tasks []ok.Task
			for _, recipe := range justfile.Recipes {
				tasks = append(tasks, ok.Task{
					Name: recipe.Name,
					RunOptions: func(ctx context.Context, args []string) (option.Options[*exec.Cmd], error) {
						return option.NewOptions(
							cob.AddArgs("--justfile", filePath),
							cob.AddArgs(recipe.Name),
							cob.AddArgs(args...),
						), nil
					},
				})
			}

			return tasks, nil
		},
	}
}
