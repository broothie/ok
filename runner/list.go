package runner

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

func (r Runner) list() error {
	filenames := make(set)
	toolNames := make(set)
	for _, task := range r.Tasks {
		filenames.insert(task.Filename())
		toolNames.insert(task.ToolName())
	}

	includeFilenames := len(filenames) > 1
	includeToolNames := len(toolNames) > 1

	lines := make([]string, len(r.Tasks))
	counter := 0
	for _, task := range r.Tasks {
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

	sort.Strings(lines)
	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
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
