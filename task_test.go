package ok

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
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

func TestOk_RunTask(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		o := &Ok{
			tasksLock: new(sync.Mutex),
			tasks: []taskInfo{
				{Task: Task{Name: "build"}, tool: &toolInfo{Tool: Tool{Name: "Make"}}},
			},
		}

		err := o.RunTask(context.Background(), "nonexistent", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nonexistent")
	})
}

func TestOk_SetUpTasks(t *testing.T) {
	t.Run("skips errored tools", func(t *testing.T) {
		o := &Ok{
			toolsLock: new(sync.Mutex),
			tasksLock: new(sync.Mutex),
			tools: []toolInfo{
				{
					Tool:           Tool{Name: "Broken"},
					commandPathErr: errors.New("not found"),
				},
			},
		}

		err := o.SetUpTasks(context.Background())
		require.NoError(t, err)
		assert.Empty(t, o.tasks)
	})

	t.Run("processes files from valid tools", func(t *testing.T) {
		o := &Ok{
			toolsLock: new(sync.Mutex),
			tasksLock: new(sync.Mutex),
			tools: []toolInfo{
				{
					Tool: Tool{
						Name: "Test",
						ProcessFile: func(ctx context.Context, filePath string, toolCfg ToolConfig) ([]Task, error) {
							return []Task{{Name: "task-from-" + filePath}}, nil
						},
					},
					commandPath: "/usr/bin/test",
					filePaths:   []string{"file1", "file2"},
				},
			},
		}

		err := o.SetUpTasks(context.Background())
		require.NoError(t, err)
		assert.Len(t, o.tasks, 2)

		names := make(map[string]bool)
		for _, tsk := range o.tasks {
			names[tsk.Name] = true
		}
		assert.True(t, names["task-from-file1"])
		assert.True(t, names["task-from-file2"])
	})

	t.Run("handles ProcessFile error", func(t *testing.T) {
		o := &Ok{
			toolsLock: new(sync.Mutex),
			tasksLock: new(sync.Mutex),
			tools: []toolInfo{
				{
					Tool: Tool{
						Name: "Failing",
						ProcessFile: func(ctx context.Context, filePath string, toolCfg ToolConfig) ([]Task, error) {
							return nil, errors.New("parse error")
						},
					},
					commandPath: "/usr/bin/failing",
					filePaths:   []string{"somefile"},
				},
			},
		}

		err := o.SetUpTasks(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse error")
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		o := &Ok{
			toolsLock: new(sync.Mutex),
			tasksLock: new(sync.Mutex),
			tools: []toolInfo{
				{
					Tool: Tool{
						Name: "Test",
						ProcessFile: func(ctx context.Context, filePath string, toolCfg ToolConfig) ([]Task, error) {
							return []Task{{Name: "should-not-run"}}, nil
						},
					},
					commandPath: "/usr/bin/test",
					filePaths:   []string{"file1"},
				},
			},
		}

		err := o.SetUpTasks(ctx)
		// Should return the context error rather than processing files
		assert.Error(t, err)
	})

	t.Run("description appears in task list", func(t *testing.T) {
		o := &Ok{
			tasks: []taskInfo{
				{
					Task:     Task{Name: "build", Description: "Build the project"},
					tool:     &toolInfo{Tool: Tool{Name: "Make"}},
					filePath: "Makefile",
				},
			},
		}

		var buf bytes.Buffer
		require.NoError(t, o.ListTasks(&buf))
		assert.Contains(t, buf.String(), "Build the project")
	})
}
