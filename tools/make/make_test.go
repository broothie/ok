package make

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tool := New()
	assert.Equal(t, "Make", tool.Name)
	assert.Equal(t, "make", tool.CommandName)
	assert.Contains(t, tool.FileGlobs, "GNUmakefile")
	assert.Contains(t, tool.FileGlobs, "Makefile")
	assert.Contains(t, tool.FileGlobs, "makefile")
	assert.NotNil(t, tool.ProcessFile)
}
