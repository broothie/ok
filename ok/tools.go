package ok

import (
	"fmt"
	"sort"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/tools"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func (ok *Ok) Registry() []tool.Tool {
	tools := funk.Filter(tools.Registry, func(t tool.Tool) bool {
		return !funk.ContainsString(ok.Options.SkipTools, t.Name())
	}).([]tool.Tool)

	sort.Slice(tools, func(i, j int) bool {
		iPriority := funk.IndexOfString(ok.Options.ToolPriority, tools[i].Name())
		if iPriority == -1 {
			return false
		}

		jPriority := funk.IndexOfString(ok.Options.ToolPriority, tools[j].Name())
		if jPriority == -1 {
			return true
		}

		return iPriority < jPriority
	})

	return funk.Reverse(tools).([]tool.Tool)
}

func (ok *Ok) List() {
	for _, tool := range ok.Registry() {
		if err := tool.Check(); err != nil {
			fmt.Printf("𝘹 %s %v\n", tool.Name(), err)
		} else {
			fmt.Printf("✔ %s\n", tool.Name())
		}
	}
}

func (ok *Ok) Init() error {
	toolName := ok.Options.Init
	tool := ok.Tool(toolName)
	if tool == nil {
		return fmt.Errorf("no tool called '%s'", toolName)
	}

	logger.Ok.Printf("initializing %s...", toolName)
	if err := tool.Init(); err != nil {
		return errors.Wrapf(err, "failed to init tool '%s'", toolName)
	}

	return nil
}

func (ok *Ok) Tool(toolName string) tool.Tool {
	return funk.Find(ok.Registry(), func(tool tool.Tool) bool { return tool.Name() == toolName }).(tool.Tool)
}
