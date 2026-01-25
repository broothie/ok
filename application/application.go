package application

import (
	"sync"

	"github.com/broothie/ok/tool"
)

type Application struct {
	Tools     []tool.Tool
	tasksLock *sync.Mutex
	Tasks     []Task
}

func New(tools []tool.Tool) *Application {
	return &Application{
		Tools:     tools,
		tasksLock: new(sync.Mutex),
	}
}
