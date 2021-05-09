package arg

import "fmt"

const help = `usage: now [options] <task> [args]

options:
	-h, --help				Print this help menu.
	-w, --watch <pattern>	Supply a glob or file pattern, and run <task> on changes. This flag can be used multiple
							times to supply multiple watch patterns.
`

func PrintHelp() {
	fmt.Print(help)
}
