package just

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tool := New()
	assert.Equal(t, "Just", tool.Name)
	assert.Equal(t, "just", tool.CommandName)
	assert.Contains(t, tool.FileGlobs, "Justfile")
	assert.Contains(t, tool.FileGlobs, "justfile")
	assert.NotNil(t, tool.ProcessFile)
}
