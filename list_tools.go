package ok

import (
	"io"
	"sort"
	"strings"

	"github.com/bobg/errors"
	"github.com/broothie/ok/table"
)

func (o *Ok) PrintTools(w io.Writer) error {
	rows := [][]string{{"TOOL", "EXECUTABLE", "FILES"}}

	for _, tl := range o.tools {
		executable := tl.CommandPath
		if tl.CommandPathErr != nil {
			executable = tl.CommandPathErr.Error()
		}

		files := strings.Join(tl.FilePaths, ",")
		if tl.FilePathsErr != nil {
			files = tl.FilePathsErr.Error()
		}

		rows = append(rows, []string{tl.Name, executable, files})
	}

	sort.Slice(rows[1:], func(i, j int) bool { return rows[1:][i][0] < rows[1:][j][0] })
	return errors.Wrap(table.Write(w, rows), "writing table")
}
