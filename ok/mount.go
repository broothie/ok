package ok

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/broothie/ok/logger"
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

type mountContext struct {
	waitGroup  *sync.WaitGroup
	errors     map[string]error
	errorsLock *sync.Mutex
}

func (ok *Ok) Mount() map[string]error {
	mountCtx := mountContext{
		waitGroup:  new(sync.WaitGroup),
		errors:     make(map[string]error),
		errorsLock: new(sync.Mutex),
	}

	tools := ok.Registry()
	tasks := make([][]Task, len(tools))
	for i, t := range tools {
		mountCtx.waitGroup.Add(1)
		go ok.mountTool(mountCtx, t, &tasks[i])
	}

	mountCtx.waitGroup.Wait()
	ok.TaskList = funk.Flatten(tasks).([]Task)
	return mountCtx.errors
}

func (ok *Ok) mountTool(ctx mountContext, tool tool.Tool, tasks *[]Task) {
	defer ctx.waitGroup.Done()

	done := make(chan struct{})
	timer := time.NewTimer(ok.Options.Timeout)

	go func() {
		defer close(done)

		start := time.Now()
		defer func() {
			if ok.Options.Debug {
				logger.Debug.Printf("mounted '%s' in %v", tool.Name(), time.Since(start))
			}
		}()

		toolConfig := tool.Config()
		if toolConfig != nil && ok.MapConfig != nil && ok.MapConfig[tool.Name()] != nil {
			if err := tomlEncodeDecode(ok.MapConfig[tool.Name()], toolConfig); err != nil {
				ctx.errorsLock.Lock()
				ctx.errors[tool.Name()] = err
				ctx.errorsLock.Unlock()
				return
			}
		}

		toolTasks, err := tool.Mount()
		if err != nil {
			ctx.errorsLock.Lock()
			ctx.errors[tool.Name()] = err
			ctx.errorsLock.Unlock()
			return
		}

		for _, toolTask := range toolTasks {
			*tasks = append(*tasks, Task{Task: toolTask, Tool: tool})
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-timer.C:
			ctx.errorsLock.Lock()
			ctx.errors[tool.Name()] = fmt.Errorf("mounting '%s' took longer than %v", tool.Name(), ok.Options.Timeout)
			ctx.errorsLock.Unlock()
			return
		}
	}
}

func tomlEncodeDecode(encode, decode interface{}) error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(encode); err != nil {
		return errors.Wrap(err, "failed to temporarily encode toml")
	}

	return errors.Wrap(toml.NewDecoder(buf).Decode(decode), "failed to decode temporary toml")
}
