package ok

import (
	"fmt"
	"io"
	"text/tabwriter"
)

func (p *Parser) WriteHelp(w io.Writer) {
	WriteVersion(w)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  $ ok [options] <task> [args]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Options:")

	t := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer t.Flush()

	for _, option := range p.availableOptions {
		var short string
		if option.Short {
			short = fmt.Sprintf("-%c", option.Name[0])
		}

		var argName string
		if option.ArgName != "" {
			argName = fmt.Sprintf(" <%s>", option.ArgName)
		}

		fmt.Fprintf(t, "\t%s\t--%s%s\t%s\n", short, option.Name, argName, option.Description)
	}
}
