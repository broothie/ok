package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/samber/lo"
)

type Task struct {
	task.Task
	Tool     tool.Tool
	Filename string
}

func (app *App) Tasks() map[string]Task {
	app.tasksOnce.Do(func() { app.tasks = app.collectTasks() })
	return app.tasks
}

func (app *App) Task(name string) (Task, bool) {
	task, found := app.Tasks()[name]
	return task, found
}

func (app *App) collectTasks() map[string]Task {
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
