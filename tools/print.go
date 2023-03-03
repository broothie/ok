package tools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
)

func (t Tools) Print() error {
	header := strings.Join([]string{"NAME", "STATUS", "EXECUTABLE"}, "\t")

	var rows []string
	for _, tool := range t {
		status := "ok"
		executable, err := exec.LookPath(tool.Config().Executable())
		if err != nil {
			status = err.Error()
		}

		rows = append(rows, strings.Join([]string{tool.Name(), status, executable}, "\t"))
	}

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