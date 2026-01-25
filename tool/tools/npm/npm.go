package npm

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok/tool"
	"github.com/broothie/option"
)

func New() tool.Tool {
	return tool.Tool{
		Name:        "NPM",
		CommandName: "npm",
		FileGlobs:   []string{"package.json"},
		ParseFile: func(ctx context.Context, filePath string) ([]tool.Task, error) {
			packageJSONFile, err := os.Open(filePath)
			if err != nil {
				return nil, errors.Wrapf(err, "opening %q", filePath)
			}
			defer packageJSONFile.Close()

			var packageJSONSchema struct {
				Scripts map[string]string `json:"scripts"`
			}

			if err := json.NewDecoder(packageJSONFile).Decode(&packageJSONSchema); err != nil {
				return nil, errors.Wrapf(err, "parsing %q", filePath)
			}

			var tasks []tool.Task
			for taskName := range packageJSONSchema.Scripts {
				tasks = append(tasks, tool.Task{
					Name: taskName,
					RunOptions: func(ctx context.Context, args []string) (option.Options[*exec.Cmd], error) {
						return option.NewOptions(
							cob.AddArgs("run", taskName),
							cob.AddArgs(args...),
						), nil
					},
				})
			}

			return tasks, nil
		},
	}
}
