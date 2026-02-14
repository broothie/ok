package nx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tool := New()
	assert.Equal(t, "Nx", tool.Name)
	assert.Equal(t, "nx", tool.CommandName)
	assert.Equal(t, []string{"nx.json"}, tool.FileGlobs)
	assert.NotNil(t, tool.ProcessFile)
}
