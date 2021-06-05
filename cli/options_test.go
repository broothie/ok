package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	taskName := "asdf"

	t.Run("no options", func(t *testing.T) {
		expected := Options{}
		_, actual, err := parserWithArgs().ParseOptions()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("task provided", func(t *testing.T) {
		expected := Options{}
		_, actual, err := parserWithArgs(taskName).ParseOptions()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("task provided with help", func(t *testing.T) {
		expected := Options{Help: true}
		_, actual, err := parserWithArgs("-h", taskName).ParseOptions()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("help", func(t *testing.T) {
		for _, option := range []string{"-h", "--help"} {
			t.Run(option, func(t *testing.T) {
				expected := Options{Help: true}
				_, actual, err := parserWithArgs(option).ParseOptions()
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			})
		}
	})

	t.Run("version", func(t *testing.T) {
		expected := Options{Version: true}
		_, actual, err := parserWithArgs("--version").ParseOptions()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("init", func(t *testing.T) {
		toolName := "asdf"
		for _, option := range []string{"-i", "--init"} {
			t.Run(option, func(t *testing.T) {
				expected := Options{Init: toolName}
				_, actual, err := parserWithArgs(option, toolName).ParseOptions()
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			})
		}
	})

	t.Run("list tools", func(t *testing.T) {
		expected := Options{ListTools: true}
		_, actual, err := parserWithArgs("--tools").ParseOptions()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("watches", func(t *testing.T) {
		t.Run("task provided with watches", func(t *testing.T) {
			expected := Options{Watches: []string{"a", "b"}}
			_, actual, err := parserWithArgs("--watch", "a", "-w", "b", taskName).ParseOptions()
			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		})

		t.Run("watches without task", func(t *testing.T) {
			_, _, err := parserWithArgs("--watch", "a", "-w", "b").ParseOptions()
			assert.EqualError(t, err, "watches provided without task")
		})

		t.Run("watches without task", func(t *testing.T) {
			_, _, err := parserWithArgs("--watch", "a", "-w", "b").ParseOptions()
			assert.EqualError(t, err, "watches provided without task")
		})
	})
}
