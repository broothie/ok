package make

import (
	"os"
	"regexp"

	"github.com/broothie/ok/stringhelp"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool/ez"
	"github.com/broothie/ok/toolhelp"
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
		Invoke: func(task ez.Task, args task.Args) *os.Process {
			return toolhelp.Exec(ToolName, task.Name()).Process
		},
	}
)
