package ok

import "sync"

type Application struct {
	toolsLock *sync.Mutex
	tools     []Tool
	tasksLock *sync.Mutex
	Tasks     []Task
}

func New() *Application {
	return &Application{
		toolsLock: new(sync.Mutex),
		tasksLock: new(sync.Mutex),
	}
}