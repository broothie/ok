package npm

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok"
	"github.com/broothie/option"
)

func New() ok.Tool {
	return ok.Tool{
		Name:        "NPM",
		CommandName: "npm",
		FileGlobs:   []string{"package.json"},
		ProcessFile: func(ctx context.Context, filePath string) ([]ok.Task, error) {
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

			var tasks []ok.Task
			for taskName := range packageJSONSchema.Scripts {
				tasks = append(tasks, ok.Task{
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
