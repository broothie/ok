package tools

import (
	"github.com/broothie/ok/tool"
	"github.com/samber/lo"
)

type Tools map[string]tool.Tool

func NewTools(tools []tool.Tool) Tools {
	return lo.Associate(tools, func(tool tool.Tool) (string, tool.Tool) { return tool.Name(), tool })
}

func FromRegistry() Tools {
	return NewTools(lo.Map(Registry(), func(newFunc tool.NewFunc, _ int) tool.Tool { return newFunc() }))
}
