package ok

import (
	"bytes"
	"sort"
	"sync"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/tools"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type Task struct {
	task.Task
	Tool tool.Tool
}

func (ok *Ok) Mount() map[string]error {
	tools := funk.Filter(tools.Registry, func(t tool.Tool) bool {
		return !funk.ContainsString(ok.Options.SkipTools, t.Name())
	}).([]tool.Tool)

	sort.Slice(tools, func(i, j int) bool {
		iPriority := funk.IndexOfString(ok.Options.SkipTools, tools[i].Name())
		if iPriority == -1 {
			return false
		}

		jPriority := funk.IndexOfString(ok.Options.SkipTools, tools[j].Name())
		if jPriority == -1 {
			return true
		}

		return iPriority < jPriority
	})

	tools = funk.Reverse(tools).([]tool.Tool)
	var mapLock sync.Mutex
	errors := make(map[string]error)
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
					errors[tool.Name()] = err
					mapLock.Unlock()
					return
				}
			}

			toolTasks, err := tool.Mount()
			if err != nil {
				mapLock.Lock()
				errors[tool.Name()] = err
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
	return errors
}

func tomlEncodeDecode(encode, decode interface{}) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(encode); err != nil {
		return errors.Wrap(err, "failed to temporarily encode toml")
	}

	return errors.Wrap(toml.NewDecoder(buf).Decode(decode), "failed to decode temporary toml")
}
