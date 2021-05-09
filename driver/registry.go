package driver

import (
	"github.com/broothie/now/task"
	"github.com/broothie/now/tool/make" // NOTE: Collides with `make` builtin
	"github.com/broothie/now/tool/ruby"
	"github.com/broothie/now/tool/yarn"
)

type MountFunc func() ([]task.Task, error)

var Registry = map[string]MountFunc{
	ruby.ToolName: ruby.Mount,
	make.ToolName: make.Mount,
	yarn.ToolName: yarn.Mount,
}
