package make

import (
	"regexp"

	"github.com/broothie/ok/stringhelp"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/toolhelp"
	"github.com/broothie/ok/tools/ez"
)

const ToolName = "make"

var (
	ruleMatcher = regexp.MustCompile(`(?m)^(?P<taskName>.*?):[^=]`)

	Make = ez.Tool{
		ToolName:             ToolName,
		CommandName:          ToolName,
		ToolFilename:         "Makefile",
		TaskMatcher:          ruleMatcher,
		CommentPrefixMatcher: stringhelp.OctothorpePrefixMatcher,
		Invoke: func(task ez.Task, args task.Args) task.RunningTask {
			return toolhelp.Exec(ToolName, task.Name())
		},
	}
)
