package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseToolOption(t *testing.T) {
	assert.Equal(t, ToolOptions{}, ParseToolOption(""))
	assert.Equal(t, ToolOptions{Name: "python"}, ParseToolOption("python"))
	assert.Equal(t, ToolOptions{Name: "python", Key: "executable"}, ParseToolOption("python.executable"))
	assert.Equal(t, ToolOptions{Name: "python", Key: "executable", Value: "/path/to/python"}, ParseToolOption("python.executable=/path/to/python"))

	assert.Equal(t, ToolOptions{Name: "docker-compose", Key: "executable", Value: "/path/to/docker"}, ParseToolOption("docker-compose.executable=/path/to/docker"))
}

func TestToolOptions_String(t *testing.T) {
	assert.Equal(t, "", ParseToolOption("").String())
	assert.Equal(t, "python", ParseToolOption("python").String())
	assert.Equal(t, "python.executable", ParseToolOption("python.executable").String())
	assert.Equal(t, "python.executable=/path/to/python", ParseToolOption("python.executable=/path/to/python").String())
}

func TestToolOptions_Action(t *testing.T) {
	assert.Equal(t, ToolOptionsActionTools, ParseToolOption("").Action())
	assert.Equal(t, ToolOptionsActionTool, ParseToolOption("python").Action())
	assert.Equal(t, ToolOptionsActionKey, ParseToolOption("python.executable").Action())
	assert.Equal(t, ToolOptionsActionSet, ParseToolOption("python.executable=/path/to/python").Action())
}
