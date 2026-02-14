package packagejson

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/broothie/ok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writePackageJSON(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "package.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestRead(t *testing.T) {
	t.Run("valid package.json with scripts", func(t *testing.T) {
		dir := t.TempDir()
		path := writePackageJSON(t, dir, `{"scripts":{"build":"tsc","test":"jest"}}`)

		s, err := read(path)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"build": "tsc", "test": "jest"}, s.Scripts)
	})

	t.Run("package.json without scripts", func(t *testing.T) {
		dir := t.TempDir()
		path := writePackageJSON(t, dir, `{"name":"my-package","version":"1.0.0"}`)

		s, err := read(path)
		require.NoError(t, err)
		assert.Nil(t, s.Scripts)
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := read("/nonexistent/package.json")
		assert.Error(t, err)
	})

	t.Run("invalid json", func(t *testing.T) {
		dir := t.TempDir()
		path := writePackageJSON(t, dir, `not json`)

		_, err := read(path)
		assert.Error(t, err)
	})
}

func TestNPM_ProcessFile(t *testing.T) {
	tool := NewNPM()
	assert.Equal(t, "NPM", tool.Name)
	assert.Equal(t, "npm", tool.CommandName)

	t.Run("discovers scripts as tasks", func(t *testing.T) {
		dir := t.TempDir()
		path := writePackageJSON(t, dir, `{"scripts":{"build":"tsc","test":"jest","lint":"eslint ."}}`)

		tasks, err := tool.ProcessFile(context.Background(), path, ok.ToolConfig{Executable: "npm"})
		require.NoError(t, err)
		require.Len(t, tasks, 3)

		taskNames := make(map[string]string)
		for _, task := range tasks {
			taskNames[task.Name] = task.Description
		}

		assert.Equal(t, "tsc", taskNames["build"])
		assert.Equal(t, "jest", taskNames["test"])
		assert.Equal(t, "eslint .", taskNames["lint"])
	})

	t.Run("no scripts returns empty", func(t *testing.T) {
		dir := t.TempDir()
		path := writePackageJSON(t, dir, `{"name":"my-package"}`)

		tasks, err := tool.ProcessFile(context.Background(), path, ok.ToolConfig{Executable: "npm"})
		require.NoError(t, err)
		assert.Empty(t, tasks)
	})
}

func TestYarn_ProcessFile(t *testing.T) {
	tool := NewYarn()
	assert.Equal(t, "Yarn", tool.Name)
	assert.Equal(t, "yarn", tool.CommandName)

	t.Run("discovers scripts as tasks", func(t *testing.T) {
		dir := t.TempDir()
		path := writePackageJSON(t, dir, `{"scripts":{"dev":"next dev","build":"next build"}}`)

		tasks, err := tool.ProcessFile(context.Background(), path, ok.ToolConfig{Executable: "yarn"})
		require.NoError(t, err)
		require.Len(t, tasks, 2)

		taskNames := make(map[string]string)
		for _, task := range tasks {
			taskNames[task.Name] = task.Description
		}

		assert.Equal(t, "next dev", taskNames["dev"])
		assert.Equal(t, "next build", taskNames["build"])
	})
}
