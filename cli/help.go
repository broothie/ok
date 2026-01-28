package cli

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/bobg/errors"
	"github.com/broothie/ok/table"
	"github.com/samber/lo"
)

var (
	//go:embed help.txt.tmpl
	helpTemplateFile string

	helpTemplate = template.Must(template.New("help").Parse(helpTemplateFile))
)

func (p *Parser) WriteHelp(w io.Writer) error {
	return errors.Wrap(helpTemplate.Execute(w, map[string]any{
		"version":   p.version,
		"flagTable": p.flagTable,
	}), "writing help")
}

func (p *Parser) flagTable() (string, error) {
	var rows [][]string
	for _, flag := range p.flags {
		shorts := ""
		if len(flag.Shorts) > 0 {
			shorts = fmt.Sprintf("-%s", string(flag.Shorts))
		}

		longNames := []string{flag.Name}
		longNames = append(longNames, flag.Aliases...)
		dashedLongs := lo.Map(longNames, func(long string, _ int) string { return fmt.Sprintf("--%s", long) })
		longs := strings.Join(dashedLongs, " ")

		rows = append(rows, []string{shorts, longs, flag.Help, fmt.Sprintf("(default: %v)", flag.DefaultValue)})
	}

	rows = table.CollapseColumns(rows)
	rows = lo.Map(rows, func(row []string, _ int) []string { return append([]string{""}, row...) })

	buffer := new(bytes.Buffer)
	if err := table.Write(buffer, rows); err != nil {
		return "", errors.Wrap(err, "writing table")
	}

	return buffer.String(), nil
}
