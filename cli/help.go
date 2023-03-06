package cli

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/broothie/ok"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

func (cli *CLI) PrintHelp(out io.Writer) error {
	strs := []string{
		fmt.Sprintf("ok %s\n", ok.Version()),
		"\n",
		"Usage:\n",
		"  ok [OPTIONS] <TASK> [TASK ARGS]\n",
		"\n",
		"Options:\n",
	}

	var table [][]string
	for _, flag := range cli.flags {
		short := ""
		if flag.hasShort() {
			short = fmt.Sprintf("-%c", flag.short)
		}

		long := fmt.Sprintf("--%s", flag.long)
		if flag.valueName != "" {
			long = fmt.Sprintf("%s %s", long, flag.valueName)
		}

		table = append(table, []string{fmt.Sprintf("  %s", short), long, flag.description})
	}

	var buf bytes.Buffer
	if err := util.PrintTable(&buf, table, 2); err != nil {
		return errors.Wrap(err, "failed to write table")
	}

	strs = append(strs, buf.String())
	if _, err := fmt.Fprint(out, strings.Join(strs, "")); err != nil {
		return errors.Wrap(err, "failed to print help")
	}

	return nil
}
