package tool

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Tools map[string]Tool

func NewTools(tools []Tool) Tools {
	return lo.Associate(tools, func(tool Tool) (string, Tool) { return tool.Name(), tool })
}

func (t Tools) Print() error {
	var rows []string
	for _, tool := range t {
		status := "ok"
		executable, err := exec.LookPath(tool.Executable())
		if err != nil {
			status = err.Error()
		}

		row := []string{tool.Name(), status, executable}
		rows = append(rows, strings.Join(row, "\t"))
	}

	sort.Strings(rows)
	rows = append([]string{strings.Join([]string{"NAME", "STATUS", "EXECUTABLE"}, "\t")}, rows...)

	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, strings.Join(rows, "\n")); err != nil {
		return errors.Wrap(err, "failed to write rows to table")
	}

	if err := table.Flush(); err != nil {
		return errors.Wrap(err, "failed to write table")
	}

	return nil
}
