package util

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func PrintTable(out io.Writer, table [][]string, padding int) error {
	rows := lo.Map(table, func(row []string, _ int) string { return strings.Join(row, "\t") })
	tableWriter := tabwriter.NewWriter(out, 0, 0, padding, ' ', 0)
	if _, err := fmt.Fprintln(tableWriter, strings.Join(rows, "\n")); err != nil {
		return errors.Wrap(err, "failed to write rows to table")
	}

	if err := tableWriter.Flush(); err != nil {
		return errors.Wrap(err, "failed to write table")
	}

	return nil
}
