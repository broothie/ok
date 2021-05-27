package tools

import (
	"sort"
	"sync"

	"github.com/broothie/ok/ok"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/thoas/go-funk"
)

func Mount(skipTools, sortTools []string) map[string]task.Task {
	ok.DebugLogger.Println(sortTools)

	tools := funk.Filter(tool.Registry, func(t tool.Tool) bool {
		return !funk.ContainsString(skipTools, t.Name())
	}).([]tool.Tool)

	sort.Slice(tools, func(i, j int) bool {
		iPriority := funk.IndexOfString(sortTools, tools[i].Name())
		if iPriority == -1 {
			return false
		}

		jPriority := funk.IndexOfString(sortTools, tools[j].Name())
		if jPriority == -1 {
			return false
		}

		return iPriority < jPriority
	})

	ok.DebugLogger.Println(funk.Map(tools, func(t tool.Tool) string { return t.Name() }).([]string))

	tasks := make(map[string]task.Task)
	var mutex sync.Mutex
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, t := range tool.Registry {
		wg.Add(1)
		go func(tool tool.Tool) {
			defer wg.Done()

			toolTasks, err := tool.Mount()
			if err != nil {
				ok.Logger.Printf("error mounting tool '%s': %v", tool.Name(), err)
				return
			}

			for _, toolTask := range toolTasks {
				name := toolTask.Name()

				mutex.Lock()
				tasks[name] = toolTask
				mutex.Unlock()
			}
		}(t)
	}

	return tasks
}
