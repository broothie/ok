package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/bobg/errors"
	"github.com/broothie/ok"
	"github.com/broothie/ok/argparser"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

const version = "v0.1.0"

type Task struct {
	ok.Task
	Tool     *ok.Tool
	FilePath string
}

func main() {
	if err := run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func run() error {
	helpFlag := &argparser.Flag{
		Name:         "help",
		Type:         argparser.FlagTypeBool,
		Help:         "Show command help.",
		Shorts:       []rune{'h'},
		DefaultValue: false,
	}

	directoryFlag := &argparser.Flag{
		Name:         "directory",
		Type:         argparser.FlagTypeString,
		Help:         "Directory to run command from.",
		Shorts:       []rune{'d'},
		DefaultValue: ".",
	}

	timeoutFlag := &argparser.Flag{
		Name:         "timeout",
		Type:         argparser.FlagTypeDuration,
		Help:         "Command timeout.",
		Shorts:       []rune{'t'},
		DefaultValue: time.Second,
	}

	argParser := argparser.New(version, os.Args[1:], helpFlag, directoryFlag, timeoutFlag)
	if err := argParser.Parse(); err != nil {
		return errors.Wrap(err, "parsing flags")
	}

	if helpFlag.Value().(bool) {
		return argParser.WriteHelp(os.Stdout)
	}

	directory := directoryFlag.Value().(string)
	timeout := timeoutFlag.Value().(time.Duration)

	if err := os.Chdir(directory); err != nil {
		return errors.Wrapf(err, "changing directory to %q", directory)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var tasks []Task
	var tasksLock sync.Mutex
	group, groupCtx := errgroup.WithContext(ctx)
	for _, tool := range ok.Registry() {
		group.Go(func() error {
			for _, glob := range tool.FileGlobs {
				if groupCtx.Err() != nil {
					return nil
				}

				filePaths, err := doublestar.Glob(os.DirFS("."), glob)
				if err != nil {
					return errors.Wrapf(err, "globbing %q for tool %q", glob, tool.Name)
				}

				for _, filePath := range filePaths {
					if groupCtx.Err() != nil {
						return nil
					}

					group.Go(func() error {
						toolTasks, err := tool.ParseFile(groupCtx, filePath)
						if err != nil {
							return errors.Wrapf(err, "parsing file %q for tool %q", filePath, tool.Name)
						}

						taskWithTools := lo.Map(toolTasks, func(task ok.Task, _ int) Task {
							return Task{
								Task:     task,
								Tool:     &tool,
								FilePath: filePath,
							}
						})

						tasksLock.Lock()
						tasks = append(tasks, taskWithTools...)
						tasksLock.Unlock()
						return nil
					})
				}
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	taskName := argParser.TaskName()
	if taskName == "" {
		if err := ListTasks(tasks); err != nil {
			return errors.Wrap(err, "listing tasks")
		}

		return nil
	}

	task, found := lo.Find(tasks, func(task Task) bool { return task.Name == taskName })
	if !found {
		return errors.Errorf("no task found with name %q", taskName)
	}

	if err := task.Run(ctx, argParser.RemainingArgs()); err != nil {
		return errors.Wrapf(err, "running task %q", task.Name)
	}

	return nil
}

func ListTasks(tasks []Task) error {
	rows := [][]string{{"TASK", "TOOL", "FILE"}}

	for _, task := range tasks {
		rows = append(rows, []string{task.Name, task.Tool.Name, task.FilePath})
	}

	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for index, row := range rows {
		if _, err := fmt.Fprintf(table, "%s\n", strings.Join(row, "\t")); err != nil {
			return errors.Wrapf(err, "writing table row %d", index)
		}
	}

	if err := table.Flush(); err != nil {
		return errors.Wrap(err, "flushing table")
	}

	return nil
}
