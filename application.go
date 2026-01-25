package ok

import "sync"

type Ok struct {
	toolsLock *sync.Mutex
	tools     []Tool
	tasksLock *sync.Mutex
	Tasks     []Task
}

func New() *Ok {
	return &Ok{
		toolsLock: new(sync.Mutex),
		tasksLock: new(sync.Mutex),
	}
}
