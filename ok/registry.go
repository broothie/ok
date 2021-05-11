package ok

import (
	"sync"

	"github.com/broothie/ok/task"
	dockercompose "github.com/broothie/ok/tool/docker-compose"
	"github.com/broothie/ok/tool/golang"
	maketool "github.com/broothie/ok/tool/make" // NOTE: Collides with `make` builtin
	"github.com/broothie/ok/tool/node"
	"github.com/broothie/ok/tool/ruby"
	"github.com/broothie/ok/tool/yarn"
)

var Registry = []Tool{
	ruby.Tool{},
	golang.Tool{},
	maketool.Tool{},
	node.Tool{},
	dockercompose.Tool{},
	yarn.Tool{},
}

func Mount() map[string]task.Task {
	tasks := make(map[string]task.Task)
	var mutex sync.Mutex
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, tool := range Registry {
		wg.Add(1)
		go func(tool Tool) {
			defer wg.Done()

			toolTasks, err := tool.Mount()
			if err != nil {
				Logger.Printf("error mounting tool '%s': %v", tool.Name(), err)
				return
			}

			for _, toolTask := range toolTasks {
				name := toolTask.Name()

				mutex.Lock()
				tasks[name] = toolTask
				mutex.Unlock()
			}
		}(tool)
	}

	return tasks
}
