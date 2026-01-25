package just

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
		Name:        "Just",
		CommandName: "just",
		FileGlobs:   []string{"justfile", "Justfile"},
		ParseFile: func(ctx context.Context, filePath string) ([]tool.Task, error) {
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

			var tasks []tool.Task
			for _, recipe := range justfile.Recipes {
				tasks = append(tasks, tool.Task{
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
