package mise

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tool := New()
	assert.Equal(t, "Mise", tool.Name)
	assert.Equal(t, "mise", tool.CommandName)
	assert.Contains(t, tool.FileGlobs, "mise.toml")
	assert.Contains(t, tool.FileGlobs, ".mise.toml")
	assert.Contains(t, tool.FileGlobs, ".mise/config.toml")
	assert.NotNil(t, tool.ProcessFile)
}
