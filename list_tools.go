package ok

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/bobg/errors"
	"github.com/broothie/ok/table"
)

func (o *Ok) ListTools(w io.Writer) error {
	rows := [][]string{{"TOOL", "EXECUTABLE", "FILES"}}

	for _, tl := range o.tools {
		executable := fmt.Sprintf("✔ %s", tl.commandPath)
		if tl.commandPathErr != nil {
			executable = fmt.Sprintf("✘ %v", tl.commandPathErr)
		}

		files := fmt.Sprintf("✔ %s", strings.Join(tl.filePaths, ","))
		if tl.filePathsErr != nil {
			files = fmt.Sprintf("✘ %v", tl.filePathsErr)
		}

		rows = append(rows, []string{tl.Name, executable, files})
	}

	sort.Slice(rows[1:], func(i, j int) bool { return rows[1:][i][0] < rows[1:][j][0] })
	return errors.Wrap(table.Write(w, rows), "writing table")
}
