package cli

import (
	"fmt"
	"testing"

	"github.com/broothie/ok/stringhelp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	taskName := "asdf"

	t.Run("no options", func(t *testing.T) {
		expected := Options{Halt: true}
		actual, err := parserWithArgs().ParseFlags()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("task provided", func(t *testing.T) {
		expected := Options{Halt: false, TaskName: taskName}
		actual, err := parserWithArgs(taskName).ParseFlags()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("task provided with help", func(t *testing.T) {
		expected := Options{Halt: true, TaskName: taskName, Help: true}
		actual, err := parserWithArgs("-h", taskName).ParseFlags()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("help", func(t *testing.T) {
		for _, flag := range []string{"-h", "--help"} {
			t.Run(flag, func(t *testing.T) {
				expected := Options{Halt: true, Help: true}
				actual, err := parserWithArgs(flag).ParseFlags()
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			})
		}
	})

	t.Run("version", func(t *testing.T) {
		expected := Options{Halt: true, Version: true}
		actual, err := parserWithArgs("--version").ParseFlags()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("init", func(t *testing.T) {
		toolName := "asdf"
		for _, flag := range []string{"-i", "--init"} {
			t.Run(flag, func(t *testing.T) {
				expected := Options{Halt: true, Init: toolName}
				actual, err := parserWithArgs(flag, toolName).ParseFlags()
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			})
		}
	})

	t.Run("list tools", func(t *testing.T) {
		expected := Options{Halt: true, ListTools: true}
		actual, err := parserWithArgs("--tools").ParseFlags()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("watches", func(t *testing.T) {
		t.Run("task provided with watches", func(t *testing.T) {
			expected := Options{Halt: false, TaskName: taskName, Watches: []string{"a", "b"}}
			actual, err := parserWithArgs("--watch", "a", "-w", "b", taskName).ParseFlags()
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})

		t.Run("watches without task", func(t *testing.T) {
			_, err := parserWithArgs("--watch", "a", "-w", "b").ParseFlags()
			assert.EqualError(t, err, "watches provided without task")
		})

		t.Run("watches without task", func(t *testing.T) {
			_, err := parserWithArgs("--watch", "a", "-w", "b").ParseFlags()
			assert.EqualError(t, err, "watches provided without task")
		})
	})

	t.Run("stoppage", func(t *testing.T) {
		argStops := map[string]bool{
			"-h":              true,
			"--help":          true,
			"--version":       true,
			"--tools":         true,
			"-w f":            false,
			"--watch f":       false,
			"-i toolName":     true,
			"--init toolName": true,
		}

		for arg, expectedStop := range argStops {
			t.Run(fmt.Sprintf("%s stops execution", arg), func(t *testing.T) {
				actual, err := parserWithArgs(append(stringhelp.SplitOnWhitespace(arg), taskName)...).ParseFlags()
				require.NoError(t, err)

				assert.Equal(t, expectedStop, actual.Halt)
			})
		}
	})
}
