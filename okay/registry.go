package okay

import (
	"sync"

	"github.com/broothie/okay/task"
	"github.com/broothie/okay/tool/golang"
	"github.com/broothie/okay/tool/make" // NOTE: Collides with `make` builtin
	"github.com/broothie/okay/tool/ruby"
	"github.com/broothie/okay/tool/yarn"
)

var Registry = map[string]Tool{
	ruby.ToolName:   ruby.Ruby{},
	golang.ToolName: golang.Golang{},
	make.ToolName:   make.Make{},
	yarn.ToolName:   yarn.Yarn{},
}

func Mount() map[string]task.Task {
	tasks := map[string]task.Task{}
	var mutex sync.Mutex
	var wg sync.WaitGroup
	defer wg.Wait()

	for toolName, t := range Registry {
		wg.Add(1)
		go func(toolName string, tool Tool) {
			defer wg.Done()

			toolTasks, err := tool.Mount()
			if err != nil {
				Logger.Printf("error mounting tool '%s': %v", toolName, err)
				return
			}

			for _, toolTask := range toolTasks {
				name := toolTask.Name()

				mutex.Lock()
				tasks[name] = toolTask
				mutex.Unlock()
			}
		}(toolName, t)
	}

	return tasks
}
