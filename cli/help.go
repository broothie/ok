package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func PrintHelp(version string) error {
	PrintVersion(version)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  $ ok [options] <task> [args]")
	fmt.Println()
	fmt.Println("Options:")

	t := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, option := range Flags {
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

		fmt.Fprintf(t, "\t%s\t--%s%s\t%s\n", short, option.Name, argName, option.Description)
	}

	return t.Flush()
}
