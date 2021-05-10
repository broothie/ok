package runner

import (
	"sync"

	"github.com/broothie/okay/tool"
)

func (r Runner) mount() {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	defer wg.Wait()

	for toolName, t := range tool.Registry {
		wg.Add(1)
		go func(toolName string, tool tool.Tool) {
			defer wg.Done()

			toolTasks, err := tool.Mount()
			if err != nil {
				log("error mounting tool '%s': %v", toolName, err)
				return
			}

			for _, toolTask := range toolTasks {
				name := toolTask.Name()

				mutex.Lock()
				r.Tasks[name] = toolTask
				mutex.Unlock()
			}
		}(toolName, t)
	}
}
