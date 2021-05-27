package tools

import (
	"sort"
	"sync"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/thoas/go-funk"
)

func Mount(skipTools, sortTools []string) task.List {
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
			return true
		}

		return iPriority < jPriority
	})

	tools = funk.Reverse(tools).([]tool.Tool)
	tasks := make([][]task.Task, len(tools))
	var wg sync.WaitGroup
	for i, t := range tools {
		wg.Add(1)
		go func(tool tool.Tool, tasks *[]task.Task) {
			defer wg.Done()

			toolTasks, err := tool.Mount()
			if err != nil {
				logger.Ok.Printf("error mounting tool '%s': %v", tool.Name(), err)
				return
			}

			for _, toolTask := range toolTasks {
				*tasks = append(*tasks, toolTask)
			}
		}(t, &tasks[i])
	}

	wg.Wait()
	return funk.Flatten(tasks).([]task.Task)
}
