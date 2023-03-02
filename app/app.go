package app

import (
	"sync"

	"github.com/broothie/ok/tool"
)

type App struct {
	Tools     tool.Tools
	tasks     Tasks
	tasksOnce *sync.Once
}

func New(tools []tool.Tool) *App {
	return &App{
		Tools:     tool.NewTools(tools),
		tasks:     make(Tasks),
		tasksOnce: new(sync.Once),
	}
}

func NewAsConfigured() *App {
	return New(Tools())
}
