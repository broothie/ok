package ok

import (
	"bytes"
	"sync"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type Task struct {
	task.Task
	Tool tool.Tool
}

func (ok *Ok) Mount() map[string]error {
	tools := ok.Registry()

	var mapLock sync.Mutex
	mountErrors := make(map[string]error)

	tasks := make([][]Task, len(tools))
	var wg sync.WaitGroup
	for i, t := range tools {
		wg.Add(1)
		go func(tool tool.Tool, tasks *[]Task) {
			defer wg.Done()

			toolConfig := tool.Config()
			if toolConfig != nil && ok.MapConfig != nil && ok.MapConfig[tool.Name()] != nil {
				if err := tomlEncodeDecode(ok.MapConfig[tool.Name()], toolConfig); err != nil {
					mapLock.Lock()
					mountErrors[tool.Name()] = err
					mapLock.Unlock()
					return
				}
			}

			toolTasks, err := tool.Mount()
			if err != nil {
				mapLock.Lock()
				mountErrors[tool.Name()] = err
				mapLock.Unlock()
				return
			}

			for _, toolTask := range toolTasks {
				*tasks = append(*tasks, Task{Task: toolTask, Tool: tool})
			}
		}(t, &tasks[i])
	}

	wg.Wait()
	ok.TaskList = funk.Flatten(tasks).([]Task)
	return mountErrors
}

func tomlEncodeDecode(encode, decode interface{}) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(encode); err != nil {
		return errors.Wrap(err, "failed to temporarily encode toml")
	}

	return errors.Wrap(toml.NewDecoder(buf).Decode(decode), "failed to decode temporary toml")
}
