package app

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/broothie/ok/parameter"
	"github.com/pkg/errors"
)

func (ok *App) ListTasks() error {
	var rows []string
	for _, task := range ok.Tasks() {
		row := []string{task.Name(), paramsString(task.Parameters()), task.Filename}
		rows = append(rows, strings.Join(row, "\t"))
	}

	sort.Strings(rows)
	rows = append([]string{strings.Join([]string{"TASK", "ARGS", "FILE"}, "\t")}, rows...)

	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, strings.Join(rows, "\n")); err != nil {
		return errors.Wrap(err, "failed to write rows to table")
	}

	if err := table.Flush(); err != nil {
		return errors.Wrap(err, "failed to write table")
	}

	return nil
}

func paramsString(params parameter.Parameters) string {
	var fields []string
	for _, param := range params {
		if param.IsRequired() {
			fields = append(fields, fmt.Sprintf("<%s>", param.Name))
		} else {
			fields = append(fields, fmt.Sprintf("--%s=%s", param.Name, *param.Default))
		}
	}

	return strings.Join(fields, " ")
}

func (ok *App) ListTools() error {
	var rows []string
	for _, tool := range ok.Tools {
		status := "ok"
		executable, err := exec.LookPath(tool.CommandName())
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
