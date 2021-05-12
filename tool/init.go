package tool

import (
	"fmt"

	"github.com/broothie/ok/ok"
	"github.com/pkg/errors"
)

func InitTool(toolName string) error {
	tool, toolExists := findTool(toolName)
	if !toolExists {
		return fmt.Errorf("no tool called '%s'", toolName)
	}

	ok.Logger.Printf("initializing %s...", toolName)
	if err := tool.Init(); err != nil {
		return errors.Wrapf(err, "failed to init tool '%s'", toolName)
	}

	return nil
}

func findTool(toolName string) (Tool, bool) {
	for _, tool := range Registry {
		if tool.Name() == toolName {
			return tool, true
		}
	}

	return nil, false
}
