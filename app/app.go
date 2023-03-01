package app

import (
	"sync"

	"github.com/broothie/ok/tool"
	"github.com/samber/lo"
)

type App struct {
	Tools map[string]tool.Tool

	tasks     map[string]Task
	tasksOnce *sync.Once
}

func New(tools []tool.Tool) *App {
	return &App{
		Tools: lo.Associate(tools, func(tool tool.Tool) (string, tool.Tool) { return tool.Name(), tool }),

		tasks:     make(map[string]Task),
		tasksOnce: new(sync.Once),
	}
}

func NewAsConfigured() *App {
	return New(Tools())
}
