package ok

import (
	"sync"
)

type toolInfo struct {
	Tool
	CommandPath    string
	CommandPathErr error
	FilePaths      []string
	FilePathsErr   error
}

type taskInfo struct {
	Task
	Tool     *toolInfo
	FilePath string
}

type Ok struct {
	toolsLock *sync.Mutex
	tools     []toolInfo
	tasksLock *sync.Mutex
	Tasks     []taskInfo
}

func New() *Ok {
	return &Ok{
		toolsLock: new(sync.Mutex),
		tasksLock: new(sync.Mutex),
	}
}
