package table

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/bobg/errors"
	"github.com/samber/lo"
)

func Write(w io.Writer, rows [][]string) error {
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

func CollapseColumns(rows [][]string) [][]string {
	columnWritten := make(map[int]bool)
	for _, row := range rows {
		for column, cell := range row {
			if cell != "" {
				columnWritten[column] = true
			}
		}
	}

	return lo.Map(rows, func(row []string, _ int) []string {
		return lo.Filter(row, func(_ string, column int) bool { return columnWritten[column] })
	})
}
