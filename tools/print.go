package tools

import (
	"io"
	"os/exec"
	"sort"

	"github.com/broothie/ok/util"
)

func (t Tools) Print(out io.Writer) error {
	table := [][]string{{"NAME", "STATUS", "EXECUTABLE"}}
	for _, tool := range t {
		status := "ok"
		executable, err := exec.LookPath(tool.Config().Executable())
		if err != nil {
			status = err.Error()
		}

		table = append(table, []string{tool.Name(), status, executable})
	}

	return util.PrintTable(out, table, 2)
}

func (t Tasks) Print(out io.Writer) error {
	headers := []string{"TASK", "ARGS", "TOOL", "FILE", "DESCRIPTION"}

	var table [][]string
	for taskName, task := range t {
		table = append(table, []string{taskName, task.Parameters().String(), task.Tool.Name(), task.Filename, task.Description()})
	}

	sort.Slice(table, func(i, j int) bool { return table[i][0] < table[j][0] })

	table = append([][]string{headers}, table...)
	return util.PrintTable(out, table, 2)
}
