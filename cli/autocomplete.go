package cli

import (
	"fmt"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tools"
	"github.com/thoas/go-funk"
)

const autocompleteOkZsh = `
#compdef ok

_ok() {
  eval "$(ok --zsh-autocomplete)"
}

_ok "$@"
`

func ZshAutocomplete() {
	builder := new(strings.Builder)

	fmt.Fprintln(builder, "local context state state_descr line")
	fmt.Fprintln(builder, "typeset -A opt_args")
	fmt.Fprintln(builder)
	fmt.Fprint(builder, "_arguments ")

	for _, flag := range Flags {
		if flag.Hidden {
			continue
		}

		description := flag.Description
		if description == "" {
			description = flag.Name
		}

		fmt.Fprint(builder, shellescape.Quote(fmt.Sprintf("--%s[%s]", flag.Name, description))+" ")

		if flag.Short {
			fmt.Fprint(builder, shellescape.Quote(fmt.Sprintf("-%c[%s]", flag.Name[0], description))+" ")
		}
	}

	builder.WriteString(shellescape.Quote(fmt.Sprintf("1:task:(%s)", strings.Join(funk.Map(tools.Mount(nil, nil), func(task task.Task) string {
		return task.Name()
	}).([]string), " "))))

	fmt.Print(builder.String())
}
