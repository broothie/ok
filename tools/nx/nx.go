package nx

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/bobg/errors"
	"github.com/broothie/cob"
	"github.com/broothie/ok"
	"github.com/broothie/option"
)

type schema struct {
	Graph graphSchema `json:"graph"`
}

type graphSchema struct {
	Nodes map[string]nodeSchema `json:"nodes"`
}

type nodeSchema struct {
	Data dataSchema `json:"data"`
}

type dataSchema struct {
	Targets map[string]any `json:"targets"`
}

func New() ok.Tool {
	return ok.Tool{
		Name:        "Nx",
		CommandName: "nx",
		FileGlobs:   []string{"nx.json"},
		ProcessFile: func(ctx context.Context, filePath string) ([]ok.Task, error) {
			// Generate the graph and output to stdout
			output, stderr, _, err := cob.Output(ctx, "nx",
				cob.AddArgs("graph"),
				cob.AddArgs("--print"),
			)
			if err != nil {
				return nil, errors.Wrapf(err, "generating nx graph: %s", stderr.String())
			}

			var payload schema
			if err := json.NewDecoder(output).Decode(&payload); err != nil {
				return nil, errors.Wrap(err, "parsing nx graph json")
			}

			var tasks []ok.Task
			for projectName, node := range payload.Graph.Nodes {
				for targetName := range node.Data.Targets {
					fullTaskName := fmt.Sprintf("%s:%s", projectName, targetName)
					tasks = append(tasks, ok.Task{
						Name: fullTaskName,
						RunOptions: func(ctx context.Context, args []string) (option.Options[*exec.Cmd], error) {
							return option.NewOptions(
								cob.AddArgs("run"),
								cob.AddArgs(fullTaskName),
								cob.AddArgs(args...),
							), nil
						},
					})
				}
			}

			return tasks, nil
		},
	}
}
