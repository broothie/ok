package ok

import (
	"sync"

	"github.com/broothie/ok/tool"
	"github.com/samber/lo"
)

type Ok struct {
	Tools map[string]tool.Tool

	tasks     map[string]Task
	tasksOnce *sync.Once
}

func New(tools []tool.Tool) *Ok {
	return &Ok{
		Tools: lo.Associate(tools, func(tool tool.Tool) (string, tool.Tool) { return tool.Name(), tool }),

		tasks:     make(map[string]Task),
		tasksOnce: new(sync.Once),
	}
}

func NewAsConfigured() *Ok {
	return New(Tools())
}
