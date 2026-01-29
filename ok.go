package ok

import (
	"sync"
)

type Ok struct {
	toolsLock *sync.Mutex
	tools     []toolInfo
	tasksLock *sync.Mutex
	tasks     []taskInfo
}

func New() *Ok {
	return &Ok{
		toolsLock: new(sync.Mutex),
		tasksLock: new(sync.Mutex),
	}
}
