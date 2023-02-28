package ok

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/tool"
	"github.com/samber/lo"
)

type Task struct {
	tool.Task
	Tool     tool.Tool
	Filename string
}

func (ok *Ok) Tasks() map[string]Task {
	ok.tasksOnce.Do(func() { ok.tasks = ok.collectTasks() })
	return ok.tasks
}

func (ok *Ok) Task(name string) (Task, bool) {
	task, found := ok.Tasks()[name]
	return task, found
}

func (ok *Ok) collectTasks() map[string]Task {
	tasks := make(map[string]Task)
	for _, tool := range ok.Tools {
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
