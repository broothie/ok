package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tool := New()
	assert.Equal(t, "Task", tool.Name)
	assert.Equal(t, "task", tool.CommandName)
	assert.Contains(t, tool.FileGlobs, "Taskfile.yml")
	assert.Contains(t, tool.FileGlobs, "Taskfile.yaml")
	assert.Contains(t, tool.FileGlobs, "taskfile.yml")
	assert.Contains(t, tool.FileGlobs, "taskfile.yaml")
	assert.NotNil(t, tool.ProcessFile)
}
