package task

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/thoas/go-funk"
)

type List []Task

func (l List) Get(taskName string) (Task, bool) {
	task := funk.Find(l, func(task Task) bool { return task.Name() == taskName })
	if task == nil {
		return nil, false
	}

	return task.(Task), true
}

func (l List) List() error {
	paramsPresent := false
	commentsPresent := false
	filenames := make(set)
	for _, task := range l {
		if task.Params().String() != "" {
			paramsPresent = true
		}

		filenames.insert(task.Filename())
		if task.Comment() != "" {
			commentsPresent = true
		}
	}

	includeFilenames := len(filenames) > 1

	lines := make([]string, len(l))
	counter := 0
	for _, task := range l {
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
