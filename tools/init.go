package tools

import (
	"fmt"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/tool"
	"github.com/pkg/errors"
)

func Init(toolName string) error {
	tool, toolExists := findTool(toolName)
	if !toolExists {
		return fmt.Errorf("no tool called '%s'", toolName)
	}

	logger.Ok.Printf("initializing %s...", toolName)
	if err := tool.Init(); err != nil {
		return errors.Wrapf(err, "failed to init tool '%s'", toolName)
	}

	return nil
}

func findTool(toolName string) (tool.Tool, bool) {
	for _, tool := range Registry {
		if tool.Name() == toolName {
			return tool, true
		}
	}

	return nil, false
}
