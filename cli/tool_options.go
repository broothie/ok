package cli

import (
	"regexp"
	"strings"

	"github.com/broothie/ok/util"
	"github.com/samber/lo"
)

type ToolOptionsAction string

const (
	ToolOptionsActionTools ToolOptionsAction = "tools"
	ToolOptionsActionTool  ToolOptionsAction = "tool"
	ToolOptionsActionKey   ToolOptionsAction = "key"
	ToolOptionsActionSet   ToolOptionsAction = "set"
)

var toolOptionParser = regexp.MustCompile(`(?P<tool>\w+)(?:\.(?P<key>\w+)(?:=(?P<value>\w+))?)?`)

type ToolOptions struct {
	Name  string
	Key   string
	Value string
}

func (t ToolOptions) String() string {
	toolKey := strings.Join(lo.Compact([]string{t.Name, t.Key}), ".")
	return strings.Join(lo.Compact([]string{toolKey, t.Value}), "=")
}

func (t ToolOptions) Action() ToolOptionsAction {
	if t.Name == "" {
		return ToolOptionsActionTools
	} else if t.Key == "" {
		return ToolOptionsActionTool
	} else if t.Value == "" {
		return ToolOptionsActionKey
	} else {
		return ToolOptionsActionSet
	}
}

// ParseToolOption valid formats:
//   python
//	 python.executable
//	 python.executable=python3
func ParseToolOption(token token) ToolOptions {
	captures := util.NamedCaptureGroups(toolOptionParser, token.String())
	return ToolOptions{
		Name:  captures["tool"],
		Key:   captures["key"],
		Value: captures["value"],
	}
}
