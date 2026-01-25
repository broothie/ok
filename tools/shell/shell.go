package shell

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/broothie/cob"
	"github.com/broothie/ok"
	"github.com/broothie/option"
)

func New() ok.Tool {
	commandName := os.Getenv("OK_SH")
	if commandName == "" {
		commandName = os.Getenv("SHELL")
	}

	if commandName == "" {
		commandName = "sh"
	}

	return ok.Tool{
		Name:        "Shell",
		CommandName: commandName,
		FileGlobs: []string{
			"**/*.bash",
			"**/*.sh",
			"**/*.zsh",
		},
		ProcessFile: func(ctx context.Context, filePath string) ([]ok.Task, error) {
			ext := filepath.Ext(filePath)

			return []ok.Task{{
				Name: strings.ReplaceAll(strings.TrimSuffix(filePath, ext), string(filepath.Separator), "."),
				RunOptions: func(ctx context.Context, args []string) (option.Options[*exec.Cmd], error) {
					return option.NewOptions(
						cob.AddArgs(filePath),
						cob.AddArgs(args...),
					), nil
				},
			}}, nil
		},
	}
}
