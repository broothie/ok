package driver

import (
	"sync"
)

func (d Driver) mount() {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	defer wg.Wait()

	for toolName, mount := range Registry {
		wg.Add(1)
		go func(toolName string, mount MountFunc) {
			defer wg.Done()

			toolTasks, err := mount()
			if err != nil {
				log("error mounting tool '%s': %v", toolName, err)
				return
			}

			for _, toolTask := range toolTasks {
				name := toolTask.Name()

				mutex.Lock()
				d.Tasks[name] = toolTask
				mutex.Unlock()
			}
		}(toolName, mount)
	}
}
