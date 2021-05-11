package okay

import (
	"fmt"
	"testing"

	"github.com/broothie/okay/tool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseOptions(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		expected := Options{Stop: true}
		actual, err := parserWithArgs().ParseOptions()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("task provided", func(t *testing.T) {
		task := "asdf"
		expected := Options{Stop: false, TaskName: task}
		actual, err := parserWithArgs(task).ParseOptions()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("task provided with help", func(t *testing.T) {
		task := "asdf"
		expected := Options{Stop: true, TaskName: task, Help: true}
		actual, err := parserWithArgs("-h", task).ParseOptions()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("help", func(t *testing.T) {
		for _, flag := range []string{"-h", "--help"} {
			t.Run(flag, func(t *testing.T) {
				expected := Options{Stop: true, Help: true}
				actual, err := parserWithArgs(flag).ParseOptions()
				assert.NoError(t, err)

				assert.Equal(t, expected, actual)
			})
		}
	})

	t.Run("version", func(t *testing.T) {
		expected := Options{Stop: true, Version: true}
		actual, err := parserWithArgs("--version").ParseOptions()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("init", func(t *testing.T) {
		toolName := "asdf"
		for _, flag := range []string{"-i", "--init"} {
			t.Run(flag, func(t *testing.T) {
				expected := Options{Stop: true, Init: toolName}
				actual, err := parserWithArgs(flag, toolName).ParseOptions()
				assert.NoError(t, err)

				assert.Equal(t, expected, actual)
			})
		}
	})

	t.Run("list tools", func(t *testing.T) {
		expected := Options{Stop: true, ListTools: true}
		actual, err := parserWithArgs("--list-tools").ParseOptions()
		assert.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("watches", func(t *testing.T) {
		t.Run("task provided with watches", func(t *testing.T) {
			task := "asdf"
			expected := Options{Stop: false, TaskName: task, Watches: []string{"a", "b"}}
			actual, err := parserWithArgs("--watch", "a", "-w", "b", task).ParseOptions()
			assert.NoError(t, err)

			assert.Equal(t, expected, actual)
		})

		t.Run("watches without task", func(t *testing.T) {
			_, err := parserWithArgs("--watch", "a", "-w", "b").ParseOptions()
			assert.EqualError(t, err, "watches provided without task")
		})

		t.Run("watches without task", func(t *testing.T) {
			_, err := parserWithArgs("--watch", "a", "-w", "b").ParseOptions()
			assert.EqualError(t, err, "watches provided without task")
		})
	})

	t.Run("stoppage", func(t *testing.T) {
		argStops := map[string]bool{
			"-h":              true,
			"--help":          true,
			"--version":       true,
			"--list-tools":    true,
			"-w f":            false,
			"--watch f":       false,
			"-i toolName":     true,
			"--init toolName": true,
		}

		for arg, expectedStop := range argStops {
			t.Run(fmt.Sprintf("%s stops execution", arg), func(t *testing.T) {
				actual, err := parserWithArgs(append(tool.SplitWhitespace(arg), "taskName")...).ParseOptions()
				require.NoError(t, err)

				assert.Equal(t, expectedStop, actual.Stop)
			})
		}
	})
}
