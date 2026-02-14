package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	all := All()
	assert.Len(t, all, 9)

	names := make(map[string]bool)
	for _, tool := range all {
		names[tool.Name] = true
		assert.NotEmpty(t, tool.Name, "tool name should not be empty")
		assert.NotEmpty(t, tool.CommandName, "command name should not be empty for %s", tool.Name)
		assert.NotEmpty(t, tool.FileGlobs, "file globs should not be empty for %s", tool.Name)
		assert.NotNil(t, tool.ProcessFile, "ProcessFile should not be nil for %s", tool.Name)
	}

	expectedTools := []string{"Just", "Make", "Mise", "Nx", "NPM", "Yarn", "Rake", "Shell", "Task"}
	for _, name := range expectedTools {
		assert.True(t, names[name], "expected tool %q to be registered", name)
	}
}
