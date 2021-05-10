package arg

import "fmt"

const help = `now v0.1.0

Usage:
  $ now [options] <task> [args]

Options:
  -h --help             Print this help menu.
     --list-tools       List all tools, and whether they're available on your machine.
  -i --init <tool>      Set up your local 
  -w --watch <pattern>  Supply a glob or file pattern, and run <task> on changes. This flag can be used multiple
                        times to supply multiple watch patterns.
`

func PrintHelp() {
	fmt.Print(help)
}
