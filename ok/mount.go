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

func (ok *Ok) Mount() map[string]error {
	mountCtx := mountContext{
		waitGroup:  new(sync.WaitGroup),
		errors:     make(map[string]error),
		errorsLock: new(sync.Mutex),
	}

	tools := funk.Reverse(ok.Registry()).([]tool.Tool)
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

	toolName := tool.Name()
	done := make(chan struct{})
	timer := time.NewTimer(ok.Options.Timeout)

	go func() {
		defer close(done)

		start := time.Now()
		defer func() {
			if ok.Options.Debug {
				logger.Debug.Printf("mounted %q in %v", toolName, time.Since(start))
			}
		}()

		if ok.MapConfig != nil && ok.MapConfig[toolName] != nil {
			if ok.Options.Debug {
				logger.Debug.Printf("%s config: %+v", toolName, ok.MapConfig[toolName])
			}

			decoder, err := tomlDecoder(ok.MapConfig[toolName])
			if err != nil {
				ctx.error(toolName, err)
			} else if err := tool.Configure(decoder); err != nil {
				ctx.error(toolName, errors.Wrapf(err, "failed to configure tool %q", toolName))
			}
		}

		toolTasks, err := tool.Mount()
		if err != nil {
			ctx.error(toolName, errors.Wrapf(err, "failed to mount tool %q", toolName))
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
			ctx.error(toolName, fmt.Errorf("mounting %q took longer than %v", toolName, ok.Options.Timeout))
			return
		}
	}
}

func tomlDecoder(encode interface{}) (*toml.Decoder, error) {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(encode); err != nil {
		return nil, errors.Wrap(err, "failed to temporarily encode toml")
	}

	return toml.NewDecoder(buf), nil
}

type mountContext struct {
	waitGroup  *sync.WaitGroup
	errors     map[string]error
	errorsLock *sync.Mutex
}

func (mc *mountContext) error(toolName string, err error) {
	mc.errorsLock.Lock()
	mc.errors[toolName] = err
	mc.errorsLock.Unlock()
}
