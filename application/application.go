package application

import (
	"sync"

	"github.com/broothie/ok/tool"
)

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
