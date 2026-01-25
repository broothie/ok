package cli

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/bobg/errors"
	"github.com/samber/lo"
)

var (
	//go:embed help.txt.tmpl
	helpTemplateFile string

	helpTemplate = template.Must(template.New("help").Parse(helpTemplateFile))
)

func (p *Parser) WriteHelp(w io.Writer) error {
	return helpTemplate.Execute(w, map[string]any{
		"version":   p.version,
		"flagTable": p.flagTable,
	})
}

func (p *Parser) flagTable() (string, error) {
	var rows []map[string]string
	shortsPresent := false
	for _, flag := range p.flags {
		shorts := ""
		if len(flag.Shorts) > 0 {
			shorts = fmt.Sprintf("-%s", string(flag.Shorts))
			shortsPresent = true
		}

		longNames := []string{flag.Name}
		longNames = append(longNames, flag.Aliases...)
		dashedLongs := lo.Map(longNames, func(long string, _ int) string { return fmt.Sprintf("--%s", long) })
		longs := strings.Join(dashedLongs, " ")

		rows = append(rows, map[string]string{
			"shorts":       shorts,
			"longs":        longs,
			"help":         flag.Help,
			"defaultValue": fmt.Sprintf("(default: %v)", flag.DefaultValue),
		})
	}

	buffer := new(bytes.Buffer)
	table := tabwriter.NewWriter(buffer, 0, 0, 2, ' ', 0)
	for _, row := range rows {
		var cells []string
		if shortsPresent {
			cells = []string{row["shorts"], row["longs"], row["help"], row["defaultValue"]}
		} else {
			cells = []string{row["longs"], row["help"], row["defaultValue"]}
		}

		if _, err := fmt.Fprintf(table, "\t%s\n", strings.Join(cells, "\t")); err != nil {
			return "", errors.Wrap(err, "writing table row")
		}
	}

	if err := table.Flush(); err != nil {
		return "", errors.Wrap(err, "writing table")
	}

	return buffer.String(), nil
}
