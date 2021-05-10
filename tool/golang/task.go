package golang

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/broothie/now/arg"
	"github.com/broothie/now/param"
	"github.com/broothie/now/task"
	"github.com/broothie/now/tool"
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

func (t Task) Invoke(args arg.Task) *os.Process {
	file, err := ioutil.TempFile("", "Nowfile-*.go")
	if err != nil {
		tool.Warn(ToolName, "failed to write go tempfile: %v", err)
		return nil
	}

	argStrings := make([]string, len(args.Positional))
	for i, positional := range args.Positional {
		parameter := t.params.PositionalRequired[i]
		switch parameter.Type {
		case param.Untyped, param.String:
			argStrings[i] = fmt.Sprintf(`"%v"`, positional)
		case param.Bool, param.Int:
			argStrings[i] = fmt.Sprint(positional)
		}
	}

	if err := tmpl.Execute(file, templateData{
		Source:   *t.fileContents,
		TaskName: t.Name(),
		Args:     argStrings,
	}); err != nil {
		tool.Warn(ToolName, "failed to write go template: %v", err)
		return nil
	}

	go file.Close()

	process := tool.Exec(ToolName, "run", file.Name()).Process
	go func() {
		process.Wait()
		os.Remove(file.Name())
	}()

	return process
}
