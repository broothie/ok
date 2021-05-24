package tool

import (
	"sync"

	"github.com/broothie/ok/ok"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tools/bash"
	dockercompose "github.com/broothie/ok/tools/docker-compose"
	"github.com/broothie/ok/tools/golang"
	maketool "github.com/broothie/ok/tools/make" // NOTE: Collides with `make` builtin
	"github.com/broothie/ok/tools/node"
	"github.com/broothie/ok/tools/python"
	"github.com/broothie/ok/tools/rake"
	"github.com/broothie/ok/tools/ruby"
	"github.com/broothie/ok/tools/yarn"
	"github.com/broothie/ok/tools/zsh"
)

var Registry = []Tool{
	bash.Bash,
	dockercompose.Tool{},
	golang.Tool{},
	maketool.Make,
	node.Tool{},
	python.Python,
	rake.Tool{},
	ruby.Ruby,
	yarn.Tool{},
	zsh.Zsh,
}

func Mount() map[string]task.Task {
	tasks := make(map[string]task.Task)
	var mutex sync.Mutex
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, t := range Registry {
		wg.Add(1)
		go func(tool Tool) {
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
