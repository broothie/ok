package task

type Base struct {
	name     string
	filename string
	toolName string
}

func NewBaseTask(name, filename, toolName string) Base {
	return Base{name: name, filename: filename, toolName: toolName}
}

func (t Base) Name() string {
	return t.name
}

func (t Base) Filename() string {
	return t.filename
}

func (t Base) ToolName() string {
	return t.toolName
}
