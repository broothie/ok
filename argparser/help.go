package argparser

import (
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/samber/lo"
)

var (
	//go:embed help.tmpl
	helpTemplateFile string

	helpTemplate = template.Must(template.New("help").Parse(helpTemplateFile))
)

func (p *ArgParser) WriteHelp(w io.Writer) error {
	return helpTemplate.Execute(w, map[string]any{
		"version": p.version,
		"flags": lo.Map(p.flags, func(f *Flag, _ int) map[string]any {
			return map[string]any{
				"shorts": fmt.Sprintf("-%s", string(f.Shorts)),
				"longs":  strings.Join(lo.Map(append([]string{f.Name}, f.Aliases...), func(long string, _ int) string { return fmt.Sprintf("--%s", long) }), " "),
				"help":   f.Help,
			}
		}),
	})
}
