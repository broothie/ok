package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/broothie/ok"
	"github.com/pkg/errors"
)

func Help() error {
	fmt.Printf("ok %s\n", ok.Version())
	fmt.Printf("\n")
	fmt.Printf("Usage:\n")
	fmt.Printf("  ok [OPTIONS] <TASK> [TASK ARGS]\n")
	fmt.Printf("\n")
	fmt.Printf("Options:\n")

	var rows []string
	for _, arg := range okArgs() {
		short := ""
		if arg.HasShort() {
			short = fmt.Sprintf("-%c", arg.Short)
		}

		row := []string{fmt.Sprintf("  %s", short), fmt.Sprintf("--%s", arg.Name), arg.Description}
		rows = append(rows, strings.Join(row, "\t"))
	}

	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, strings.Join(rows, "\n")); err != nil {
		return errors.Wrap(err, "failed to write rows to table")
	}

	if err := table.Flush(); err != nil {
		return errors.Wrap(err, "failed to write table")
	}

	return nil
}
