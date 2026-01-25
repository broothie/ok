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
		// "flags": lo.Map(p.flags, func(f *Flag, _ int) map[string]any {
		// 	return map[string]any{
		// 		"shorts": fmt.Sprintf("-%s", string(f.Shorts)),
		// 		"longs":  strings.Join(lo.Map(append([]string{f.Name}, f.Aliases...), func(long string, _ int) string { return fmt.Sprintf("--%s", long) }), " "),
		// 		"help":   f.Help,
		// 	}
		// }),
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

		if _, err := fmt.Fprintln(table, fmt.Sprintf("\t%s", strings.Join(cells, "\t"))); err != nil {
			return "", errors.Wrap(err, "")
		}
	}

	if err := table.Flush(); err != nil {
		return "", errors.Wrap(err, "writing table")
	}

	return buffer.String(), nil
}
