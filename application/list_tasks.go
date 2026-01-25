package application

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/bobg/errors"
)

func (a *Application) ListTasks(w io.Writer) error {
	rows := [][]string{{"TASK", "TOOL", "FILE"}}

	for _, task := range a.Tasks {
		rows = append(rows, []string{task.Name, task.Tool.Name, task.FilePath})
	}

	table := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for index, row := range rows {
		if _, err := fmt.Fprintf(table, "%s\n", strings.Join(row, "\t")); err != nil {
			return errors.Wrapf(err, "writing table row %d", index)
		}
	}

	if err := table.Flush(); err != nil {
		return errors.Wrap(err, "flushing table")
	}

	return nil
}
