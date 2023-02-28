package main

import (
	"context"
	"fmt"

	"github.com/broothie/ok"
	"github.com/broothie/ok/cli"
	"github.com/broothie/ok/logger"
)

func main() {
	parser := cli.NewFromArgs()
	okArgs, err := parser.ParseOkArgs()
	if err != nil {
		logger.Log.Fatalf("failed to parse args: %v", err)
	}

	if okArgs.Help {
		fmt.Println("helping...")
		return
	}

	ok := ok.NewAsConfigured()
	if okArgs.ListTools {
		if err := ok.ListTools(); err != nil {
			logger.Log.Fatalf("failed to list tools: %v", err)
		}

		return
	}

	if okArgs.TaskName == "" {
		if err := ok.ListTasks(); err != nil {
			logger.Log.Fatalf("failed to list tasks: %v", err)
		}

		return
	}

	task, found := ok.Task(okArgs.TaskName)
	if !found {
		logger.Log.Fatalf("no task found with name %q", okArgs.TaskName)
	}

	args, err := parser.ParseTaskArgs(task.Parameters())
	if err != nil {
		logger.Log.Fatalf("failed to parse task args: %v", err)
	}

	if err := task.Run(context.Background(), args); err != nil {
		logger.Log.Fatalf("failed to run task %q: %v", okArgs.TaskName, err)
	}
}
