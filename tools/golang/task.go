package golang

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/broothie/okay/arg"
	"github.com/broothie/okay/param"
	"github.com/broothie/okay/task"
	"github.com/broothie/okay/toolhelp"
)

var tmpl = template.Must(template.New("").Parse(`{{ .Source }}

func main() {
	{{ .TaskName }}(
		{{ range $arg := .Args }}
			{{ $arg }}, 
		{{ end }}
	)
}
`))

type templateData struct {
	Source   string
	TaskName string
	Args     []string
}

type Task struct {
	task.Base
	params       param.Params
	fileContents *string
}

func (t Task) Params() param.Params {
	return t.params
}

func (t Task) Invoke(args arg.Args) *os.Process {
	file, err := ioutil.TempFile("", "Okayfile-*.go")
	if err != nil {
		toolhelp.Warn(ToolName, "failed to write go tempfile: %v", err)
		return nil
	}

	argStrings := make([]string, len(args.Positional))
	for i, positional := range args.Positional {
		//parameter := t.params.PositionalRequired[i]
		switch positional.Param.Type {
		case param.Untyped, param.String:
			argStrings[i] = fmt.Sprintf(`"%v"`, positional.Value)
		case param.Bool, param.Int:
			argStrings[i] = fmt.Sprint(positional.Value)
		}
	}

	if err := tmpl.Execute(file, templateData{
		Source:   *t.fileContents,
		TaskName: t.Name(),
		Args:     argStrings,
	}); err != nil {
		toolhelp.Warn(ToolName, "failed to write go template: %v", err)
		return nil
	}

	go file.Close()

	process := toolhelp.Exec(ToolName, "run", file.Name()).Process
	go func() {
		process.Wait()
		os.Remove(file.Name())
	}()

	return process
}
