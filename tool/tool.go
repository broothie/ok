package tool

import (
	"github.com/broothie/ok/task"
	"github.com/pelletier/go-toml"
)

type Tool interface {
	Name() string
	Init() error
	Check() error
	Configure(*toml.Decoder) error
	Mount() ([]task.Task, error)
}
