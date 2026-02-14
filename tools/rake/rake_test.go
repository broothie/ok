package rake

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitRakeArgs(t *testing.T) {
	t.Run("empty args", func(t *testing.T) {
		taskArgs, rakeVars, err := splitRakeArgs(nil)
		require.NoError(t, err)
		assert.Nil(t, taskArgs)
		assert.Nil(t, rakeVars)
	})

	t.Run("task args only", func(t *testing.T) {
		taskArgs, rakeVars, err := splitRakeArgs([]string{"arg1", "arg2"})
		require.NoError(t, err)
		assert.Equal(t, []string{"arg1", "arg2"}, taskArgs)
		assert.Nil(t, rakeVars)
	})

	t.Run("rake vars only", func(t *testing.T) {
		taskArgs, rakeVars, err := splitRakeArgs([]string{"FOO=bar", "BAZ=qux"})
		require.NoError(t, err)
		assert.Nil(t, taskArgs)
		assert.Equal(t, []string{"FOO=bar", "BAZ=qux"}, rakeVars)
	})

	t.Run("mixed args and vars", func(t *testing.T) {
		taskArgs, rakeVars, err := splitRakeArgs([]string{"arg1", "FOO=bar", "arg2", "BAZ=qux"})
		require.NoError(t, err)
		assert.Equal(t, []string{"arg1", "arg2"}, taskArgs)
		assert.Equal(t, []string{"FOO=bar", "BAZ=qux"}, rakeVars)
	})

	t.Run("flags treated as task args", func(t *testing.T) {
		taskArgs, rakeVars, err := splitRakeArgs([]string{"--verbose", "-t"})
		require.NoError(t, err)
		assert.Equal(t, []string{"--verbose", "-t"}, taskArgs)
		assert.Nil(t, rakeVars)
	})

	t.Run("flag with equals is not a rake var", func(t *testing.T) {
		taskArgs, rakeVars, err := splitRakeArgs([]string{"--config=value"})
		require.NoError(t, err)
		assert.Equal(t, []string{"--config=value"}, taskArgs)
		assert.Nil(t, rakeVars)
	})

	t.Run("error on comma in arg", func(t *testing.T) {
		_, _, err := splitRakeArgs([]string{"arg,with,commas"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "','")
	})

	t.Run("error on bracket in arg", func(t *testing.T) {
		_, _, err := splitRakeArgs([]string{"arg]bracket"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "']'")
	})
}

func TestNew(t *testing.T) {
	tool := New()
	assert.Equal(t, "Rake", tool.Name)
	assert.Equal(t, "rake", tool.CommandName)
	assert.Contains(t, tool.FileGlobs, "Rakefile")
	assert.Contains(t, tool.FileGlobs, "rakefile")
	assert.Contains(t, tool.FileGlobs, "Rakefile.rb")
	assert.Contains(t, tool.FileGlobs, "rakefile.rb")
}
