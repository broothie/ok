package ok

import (
	"sort"
	"sync"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/tools"
	"github.com/thoas/go-funk"
)

type Task struct {
	task.Task
	Tool tool.Tool
}

func (ok *Ok) Mount() map[string]error {
	filteredTools := funk.Filter(tools.Registry, func(t tool.Tool) bool {
		return !funk.ContainsString(ok.Options.SkipTools, t.Name())
	}).([]tool.Tool)

	sort.Slice(filteredTools, func(i, j int) bool {
		iPriority := funk.IndexOfString(ok.Options.SkipTools, filteredTools[i].Name())
		if iPriority == -1 {
			return false
		}

		jPriority := funk.IndexOfString(ok.Options.SkipTools, filteredTools[j].Name())
		if jPriority == -1 {
			return true
		}

		return iPriority < jPriority
	})

	var mapLock sync.Mutex
	errors := make(map[string]error)
	reversedTools := funk.Reverse(filteredTools).([]tool.Tool)
	tasks := make([][]Task, len(tools.Registry))
	var wg sync.WaitGroup
	for i, t := range reversedTools {
		wg.Add(1)
		go func(tool tool.Tool, tasks *[]Task) {
			defer wg.Done()

			toolTasks, err := tool.Mount()
			if err != nil {
				mapLock.Lock()
				errors[tool.Name()] = err
				mapLock.Unlock()
				return
			}

			for _, toolTask := range toolTasks {
				*tasks = append(*tasks, Task{
					Task: toolTask,
					Tool: tool,
				})
			}
		}(t, &tasks[i])
	}

	wg.Wait()
	ok.TaskList = funk.Flatten(tasks).([]Task)
	return errors
}
