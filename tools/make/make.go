package make

import (
	"regexp"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tools/ez"
	"github.com/broothie/ok/util"
)

const ToolName = "make"

var (
	ruleMatcher = regexp.MustCompile(`(?m)^(?P<taskName>[\w-]*?):[^=]?`)

	Make = ez.Tool{
		ToolName:             ToolName,
		CommandName:          ToolName,
		ToolFilename:         "Makefile",
		TaskMatcher:          ruleMatcher,
		CommentPrefixMatcher: util.OctothorpePrefixMatcher,
		Invoke: func(task ez.Task, args task.Args) task.RunningTask {
			return util.Exec(ToolName, task.Name())
		},
	}
)
