package ok

import (
	"context"
	"encoding/json"
	"os"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
)

func NewNPM() Tool {
	return Tool{
		Name:        "NPM",
		CommandName: "npm",
		FileGlobs:   []string{"package.json"},
		ParseFile: func(ctx context.Context, filePath string) ([]Task, error) {
			packageJSONFile, err := os.Open(filePath)
			if err != nil {
				return nil, errors.Wrapf(err, "opening %q", filePath)
			}

			var packageJSONSchema struct {
				Scripts map[string]string `json:"scripts"`
			}

			if err := json.NewDecoder(packageJSONFile).Decode(&packageJSONSchema); err != nil {
				return nil, errors.Wrapf(err, "parsing %q", filePath)
			}

			var tasks []Task
			for taskName := range packageJSONSchema.Scripts {
				tasks = append(tasks, Task{
					Name: taskName,
					Run: func(ctx context.Context, args []string) error {
						_, err := cob.Run(ctx, "npm",
							cob.AddArgs("run", taskName),
							cob.AddArgs(args...),
							cob.SetStdin(os.Stdin),
							cob.SetStdout(os.Stdout),
							cob.SetStderr(os.Stderr),
						)

						return err
					},
				})
			}

			return tasks, nil
		},
	}
}
