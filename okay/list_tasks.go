package okay

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/broothie/okay/task"
)

func ListTasks(w io.Writer, tasks map[string]task.Task) error {
	filenames := make(set)
	toolNames := make(set)
	for _, task := range tasks {
		filenames.insert(task.Filename())
		toolNames.insert(task.ToolName())
	}

	includeFilenames := len(filenames) > 1
	includeToolNames := len(toolNames) > 1

	lines := make([]string, len(tasks))
	counter := 0
	for _, task := range tasks {
		columns := []string{fmt.Sprintf("%s %s", task.Name(), task.Params())}

		if includeFilenames {
			columns = append(columns, task.Filename())
		}

		if includeToolNames {
			columns = append(columns, task.ToolName())
		}

		lines[counter] = fmt.Sprintf("%s\n", strings.Join(columns, "\t"))
		counter++
	}

	headerSlice := []string{"TASK"}
	if includeFilenames {
		headerSlice = append(headerSlice, "FILE")
	}

	if includeToolNames {
		headerSlice = append(headerSlice, "TOOL")
	}

	sort.Strings(lines)
	lines = append([]string{strings.Join(headerSlice, "\t") + "\n"}, lines...)

	table := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, line := range lines {
		if _, err := fmt.Fprint(table, line); err != nil {
			return err
		}
	}

	return table.Flush()
}

type set map[interface{}]struct{}

var present struct{}

func (s set) insert(v interface{}) {
	s[v] = present
}
