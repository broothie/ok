package tools

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tool"
	"github.com/pkg/errors"
)

type Task struct {
	task.Task
	Tool     tool.Tool
	Filename string
}

type Tasks map[string]Task

func (t Tasks) Task(name string) (Task, bool) {
	task, found := t[name]
	return task, found
}

func (t Tasks) Print() error {
	header := strings.Join([]string{"TASK", "ARGS", "TOOL", "FILE", "DESCRIPTION"}, "\t")

	var rows []string
	for taskName, task := range t {
		row := []string{taskName, task.Parameters().String(), task.Tool.Name(), task.Filename, task.Description()}
		rows = append(rows, strings.Join(row, "\t"))
	}

	sort.Strings(rows)
	rows = append([]string{header}, rows...)

	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, strings.Join(rows, "\n")); err != nil {
		return errors.Wrap(err, "failed to write rows to table")
	}

	if err := table.Flush(); err != nil {
		return errors.Wrap(err, "failed to write table")
	}

	return nil
}
