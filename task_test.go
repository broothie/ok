package ok

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOk_ListTasks(t *testing.T) {
	type testCase struct {
		tasks            []taskInfo
		expectedContains []string
		expectedOrder    []string
	}

	makeTool := &toolInfo{Tool: Tool{Name: "Make"}}
	rakeTool := &toolInfo{Tool: Tool{Name: "Rake"}}

	tests := map[string]testCase{
		"no tasks": {
			tasks:            nil,
			expectedContains: []string{"TASK", "TOOL", "FILE"},
			expectedOrder:    nil,
		},
		"single task": {
			tasks: []taskInfo{
				{
					Task:     Task{Name: "build"},
					tool:     makeTool,
					filePath: "Makefile",
				},
			},
			expectedContains: []string{"build", "Make", "Makefile"},
			expectedOrder:    []string{"build"},
		},
		"multiple tasks sorted alphabetically": {
			tasks: []taskInfo{
				{
					Task:     Task{Name: "test"},
					tool:     rakeTool,
					filePath: "Rakefile",
				},
				{
					Task:     Task{Name: "build"},
					tool:     makeTool,
					filePath: "Makefile",
				},
			},
			expectedContains: []string{"build", "test", "Make", "Rake", "Makefile", "Rakefile"},
			expectedOrder:    []string{"build", "test"},
		},
		"tasks from same tool": {
			tasks: []taskInfo{
				{
					Task:     Task{Name: "build"},
					tool:     makeTool,
					filePath: "Makefile",
				},
				{
					Task:     Task{Name: "clean"},
					tool:     makeTool,
					filePath: "Makefile",
				},
			},
			expectedContains: []string{"build", "clean", "Make", "Makefile"},
			expectedOrder:    []string{"build", "clean"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			o := &Ok{tasks: tc.tasks}

			var buf bytes.Buffer
			err := o.ListTasks(&buf)

			require.NoError(t, err)

			output := buf.String()
			for _, s := range tc.expectedContains {
				assert.Contains(t, output, s)
			}

			// Verify sort order
			if len(tc.expectedOrder) > 1 {
				for i := 0; i < len(tc.expectedOrder)-1; i++ {
					pos1 := strings.Index(output, tc.expectedOrder[i])
					pos2 := strings.Index(output, tc.expectedOrder[i+1])
					assert.Less(t, pos1, pos2, "expected %q before %q", tc.expectedOrder[i], tc.expectedOrder[i+1])
				}
			}
		})
	}
}
