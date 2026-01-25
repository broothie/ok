package application

import (
	"sync"

	"github.com/broothie/ok/tool"
)

type Task struct {
	tool.Task
	Tool     *tool.Tool
	FilePath string
}

type Application struct {
	workingDirectory string
	Tools            []tool.Tool
	tasksLock        *sync.Mutex
	Tasks            []Task
}

func New(workingDirectory string, tools []tool.Tool) *Application {
	return &Application{
		workingDirectory: workingDirectory,
		Tools:            tools,
		tasksLock:        new(sync.Mutex),
	}
}
