package ok

import (
	"io"
	"sort"

	"github.com/bobg/errors"
	"github.com/broothie/ok/table"
)

func (o *Ok) ListTasks(w io.Writer) error {
	rows := [][]string{{"TASK", "TOOL", "FILE"}}

	for _, task := range o.Tasks {
		rows = append(rows, []string{task.Name, task.Tool.Name, task.FilePath})
	}

	sort.Slice(rows[1:], func(i, j int) bool { return rows[1:][i][0] < rows[1:][j][0] })
	return errors.Wrap(table.Write(w, rows), "writing table")
}
