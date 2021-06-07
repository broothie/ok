package ok

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

func (ok *Ok) ListTasks() error {
	paramsPresent := false
	commentsPresent := false
	filenames := make(set)
	for _, task := range ok.TaskList {
		if task.Params().String() != "" {
			paramsPresent = true
		}

		filenames.insert(task.Filename())
		if task.Comment() != "" {
			commentsPresent = true
		}
	}

	includeFilenames := len(filenames) > 1

	lines := make([]string, len(ok.TaskList))
	counter := 0
	for _, task := range ok.TaskList {
		columns := []string{task.Name()}
		if paramsPresent {
			columns = append(columns, task.Params().String())
		}

		if includeFilenames {
			columns = append(columns, task.Filename())
		}

		if commentsPresent {
			columns = append(columns, task.Comment())
		}

		lines[counter] = fmt.Sprintf("%s\n", strings.Join(columns, "\t"))
		counter++
	}

	headerSlice := []string{"TASK"}
	if paramsPresent {
		headerSlice = append(headerSlice, "ARGS")
	}

	if includeFilenames {
		headerSlice = append(headerSlice, "FILE")
	}

	if commentsPresent {
		headerSlice = append(headerSlice, "DESCRIPTION")
	}

	sort.Strings(lines)
	if len(headerSlice) > 1 {
		lines = append([]string{strings.Join(headerSlice, "\t") + "\n"}, lines...)
	}

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
