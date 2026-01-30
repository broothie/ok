package just

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
	Recipes map[string]recipeSchema `json:"recipes"`
}

type recipeSchema struct {
	Name string  `json:"name"`
	Doc  *string `json:"doc"`
}

func New() ok.Tool {
	return ok.Tool{
		Name:        "Just",
		CommandName: "just",
		FileGlobs:   []string{"Justfile", "justfile"},
		ProcessFile: func(ctx context.Context, filePath string, toolCfg ok.ToolConfig) ([]ok.Task, error) {
			output, _, _, err := cob.Output(ctx, toolCfg.Executable,
				cob.AddArgs("--justfile", filePath),
				cob.AddArgs("--dump"),
				cob.AddArgs("--dump-format", "json"),
			)
			if err != nil {
				return nil, errors.Wrapf(err, "listing just recipes from %q", filePath)
			}

			var payload schema
			if err := json.NewDecoder(output).Decode(&payload); err != nil {
				return nil, errors.Wrapf(err, "parsing just recipes from %q", filePath)
			}

			return lo.Map(lo.Values(payload.Recipes), func(recipe recipeSchema, _ int) ok.Task {
				var description string
				if recipe.Doc != nil {
					description = *recipe.Doc
				}

				return ok.Task{
					Name:        recipe.Name,
					Description: description,
					RunOptions: func(ctx context.Context, args []string, toolCfg ok.ToolConfig) (option.Options[*exec.Cmd], error) {
						return option.NewOptions(
							cob.AddArgs("--justfile", filePath),
							cob.AddArgs(recipe.Name),
							cob.AddArgs(args...),
						), nil
					},
				}
			}), nil
		},
	}
}
