package main

import (
	"fmt"
	"os"

	"github.com/broothie/okay/okay"
)

func main() {
	// Parse options
	parser := okay.NewParser(os.Args[1:])
	options, err := parser.ParseOptions()
	if err != nil {
		okay.Logger.Println(err)
		os.Exit(1)
		return
	}

	// Process options
	switch {
	case options.Help:
		parser.WriteHelp(os.Stdout)

	case options.Version:
		fmt.Printf("ok v%s", okay.Version)

	case options.Init != "":
		if err := okay.InitTool(options.Init); err != nil {
			okay.Logger.Println(err)
			os.Exit(1)
			return
		}

	case options.ListTools:
		okay.ListTools(os.Stdout)

	case options.TaskName == "":
		options.Stop = true
		if err := okay.ListTasks(os.Stdout, okay.Mount()); err != nil {
			okay.Logger.Println(err)
			os.Exit(1)
			return
		}
	}

	if options.Stop {
		return
	}

	// Get task
	taskName := options.TaskName
	tasks := okay.Mount()
	task, taskExists := tasks[options.TaskName]
	if !taskExists {
		okay.Logger.Printf("no task called '%s'", taskName)
		os.Exit(1)
		return
	}

	// Parse task args
	args, err := parser.ParseArgs(task.Params())
	if err != nil {
		okay.Logger.Println(err)
		os.Exit(1)
		return
	}

	// Run task
	if _, err := task.Invoke(args).Wait(); err != nil {
		okay.Logger.Println(err)
		os.Exit(1)
	}
}
