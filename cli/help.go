package cli

import (
	"bufio"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/pkg/errors"
)

const header = `
Usage:
  $ ok [options] <task> [args]

Options:
`

func PrintHelp(w io.Writer, version string) error {
	buf := bufio.NewWriter(w)
	if err := PrintVersion(buf, version); err != nil {
		return err
	}

	if _, err := fmt.Fprint(buf, header); err != nil {
		return errors.Wrap(err, "failed to write help line")
	}

	table := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	for _, option := range options {
		if option.Hidden {
			continue
		}

		var short string
		if option.Short {
			short = fmt.Sprintf("-%c", option.Name[0])
		}

		var argName string
		if option.ArgName != "" {
			argName = fmt.Sprintf(" <%s>", option.ArgName)
		}

		if _, err := fmt.Fprintf(table, "\t%s\t--%s%s\t%s\n", short, option.Name, argName, option.Description); err != nil {
			return errors.Wrap(err, "failed to write help line")
		}
	}

	if err := table.Flush(); err != nil {
		return errors.Wrap(err, "failed to flush help table")
	}

	return errors.Wrap(buf.Flush(), "failed to print help")
}
