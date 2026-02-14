package shell

import (
	"context"
	"testing"

	"github.com/broothie/ok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tool := New()
	assert.Equal(t, "Shell", tool.Name)
	assert.Contains(t, tool.FileGlobs, "**/*.sh")
	assert.Contains(t, tool.FileGlobs, "**/*.bash")
	assert.Contains(t, tool.FileGlobs, "**/*.zsh")
}

func TestProcessFile_TaskName(t *testing.T) {
	tests := map[string]struct {
		filePath     string
		expectedName string
	}{
		"simple sh file":     {"build.sh", "build"},
		"simple bash file":   {"deploy.bash", "deploy"},
		"simple zsh file":    {"setup.zsh", "setup"},
		"nested path":        {"scripts/build.sh", "scripts.build"},
		"deeply nested path": {"scripts/ci/test.sh", "scripts.ci.test"},
	}

	tool := New()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tasks, err := tool.ProcessFile(context.Background(), tc.filePath, ok.ToolConfig{})
			require.NoError(t, err)
			require.Len(t, tasks, 1)
			assert.Equal(t, tc.expectedName, tasks[0].Name)
		})
	}
}
