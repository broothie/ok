package task

import "github.com/pelletier/go-toml"

type Base struct {
	name     string
	filename string
	toolName string
}

func NewBase(name, filename, toolName string) Base {
	return Base{name: name, filename: filename, toolName: toolName}
}

func (b Base) Name() string {
	return b.name
}

func (b Base) Comment() string {
	return ""
}

func (b Base) Filename() string {
	return b.filename
}

func (b Base) ToolName() string {
	return b.toolName
}

func (b Base) Configure(*toml.Decoder) error {
	return nil
}

func (Base) Params() Parameters {
	return Parameters{}
}
