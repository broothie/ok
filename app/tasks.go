package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Task struct {
	task.Task
	Tool     tool.Tool
	Filename string
}

type Tasks map[string]Task

func (t Tasks) Task(name string) (Task, bool) {
	task, found := t[name]
	return task, found
}

func (t Tasks) Print() error {
	var rows []string
	for _, task := range t {
		row := []string{task.Name(), task.Parameters().String(), task.Filename}
		rows = append(rows, strings.Join(row, "\t"))
	}

	sort.Strings(rows)
	header := strings.Join([]string{"TASK", "ARGS", "FILE"}, "\t")
	rows = append([]string{header}, rows...)

	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, strings.Join(rows, "\n")); err != nil {
		return errors.Wrap(err, "failed to write rows to table")
	}

	if err := table.Flush(); err != nil {
		return errors.Wrap(err, "failed to write table")
	}

	return nil
}

func (app *App) Tasks() Tasks {
	app.tasksOnce.Do(func() { app.tasks = app.collectTasks() })
	return app.tasks
}

func (app *App) collectTasks() Tasks {
	tasks := make(map[string]Task)
	for _, tool := range app.Tools {
		paths := append(tool.Filenames(), lo.FlatMap(tool.Extensions(), func(extension string, _ int) []string {
			okfiles := []string{fmt.Sprintf("Okfile.%s", extension)}
			subOkFiles, err := filepath.Glob(fmt.Sprintf("okfiles/*.%s", extension))
			if err != nil {
				logger.Log.Fatal("failed to glob for extension", extension)
			}

			return append(okfiles, subOkFiles...)
		})...)

		for _, filename := range paths {
			if _, err := os.Stat(filename); err != nil {
				if os.IsNotExist(err) {
					continue
				}

				logger.Log.Fatal("failed to stat file", filename, err)
			}

			toolTasks, err := tool.ProcessFile(filename)
			if err != nil {
				logger.Log.Fatal("failed to parse file", filename, err)
			}

			for _, task := range toolTasks {
				tasks[task.Name()] = Task{Task: task, Tool: tool, Filename: filename}
			}
		}
	}

	return tasks
}
