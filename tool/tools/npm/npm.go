package npm

import (
	"context"
	"encoding/json"
	"os"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok/tool"
)

func NewNPM() tool.Tool {
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
