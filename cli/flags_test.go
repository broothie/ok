package cli

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCLI(args ...string) *CLI {
	return &CLI{
		parser:      newParser(args),
		flags:       flags(),
		optionsOnce: new(sync.Once),
	}
}

func Test_flags(t *testing.T) {
	t.Run("-h", func(t *testing.T) {
		cli := New([]string{"-h"})
		options, err := cli.Options()
		require.NoError(t, err)

		assert.True(t, options.Help)
	})

	t.Run("--help", func(t *testing.T) {
		cli := New([]string{"--help"})
		options, err := cli.Options()
		require.NoError(t, err)

		assert.True(t, options.Help)
	})

	t.Run("-V", func(t *testing.T) {
		cli := New([]string{"-V"})
		options, err := cli.Options()
		require.NoError(t, err)

		assert.True(t, options.Version)
	})

	t.Run("--version", func(t *testing.T) {
		cli := New([]string{"--version"})
		options, err := cli.Options()
		require.NoError(t, err)

		assert.True(t, options.Version)
	})

	t.Run("--tools", func(t *testing.T) {
		cli := New([]string{"--tools"})
		options, err := cli.Options()
		require.NoError(t, err)

		assert.True(t, options.ListTools)
	})

	t.Run("--tool", func(t *testing.T) {
		t.Run("", func(t *testing.T) {
			cli := New([]string{"--tool"})
			options, err := cli.Options()
			require.NoError(t, err)

			assert.Equal(t, []ToolOptions{{}}, options.ToolOptions)
		})

		t.Run("python", func(t *testing.T) {
			cli := New([]string{"--tool", "python"})
			options, err := cli.Options()
			require.NoError(t, err)

			assert.Equal(t, []ToolOptions{{Name: "python"}}, options.ToolOptions)
		})

		t.Run("python.executable", func(t *testing.T) {
			cli := New([]string{"--tool", "python.executable"})
			options, err := cli.Options()
			require.NoError(t, err)

			assert.Equal(t, []ToolOptions{{Name: "python", Key: "executable"}}, options.ToolOptions)
		})

		t.Run("python.executable=/path/to/python", func(t *testing.T) {
			cli := New([]string{"--tool", "python.executable=/path/to/python"})
			options, err := cli.Options()
			require.NoError(t, err)

			assert.Equal(t, []ToolOptions{{Name: "python", Key: "executable", Value: "/path/to/python"}}, options.ToolOptions)
		})
	})

	t.Run("--init", func(t *testing.T) {
		cli := New([]string{"--init", "python"})
		options, err := cli.Options()
		require.NoError(t, err)

		assert.Equal(t, options.InitTool, "python")
	})

	t.Run("-w somefile", func(t *testing.T) {
		cli := New([]string{"-w", "somefile"})
		options, err := cli.Options()
		require.NoError(t, err)

		assert.Equal(t, options.Watches, []string{"somefile"})
	})

	t.Run("--watch somefile", func(t *testing.T) {
		cli := New([]string{"--watch", "somefile"})
		options, err := cli.Options()
		require.NoError(t, err)

		assert.Equal(t, options.Watches, []string{"somefile"})
	})

	t.Run("-w somefile --watch otherfile", func(t *testing.T) {
		cli := New([]string{"-w", "somefile", "--watch", "otherfile"})
		options, err := cli.Options()
		require.NoError(t, err)

		assert.Equal(t, options.Watches, []string{"somefile", "otherfile"})
	})
}
