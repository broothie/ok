package tools

import (
	"sync"

	"github.com/broothie/ok/ok"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
)

func Mount() map[string]task.Task {
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
