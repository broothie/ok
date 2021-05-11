package golang

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/broothie/okay/tool"

	"github.com/broothie/okay/task"
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
	params       task.Parameters
	fileContents *string
}

func (t Task) Params() task.Parameters {
	return t.params
}

func (t Task) Invoke(args task.Args) *os.Process {
	file, err := ioutil.TempFile("", "Okayfile-*.go")
	if err != nil {
		tool.Warn(ToolName, "failed to write go tempfile: %v", err)
		return nil
	}

	argStrings := make([]string, len(args.Positional))
	for i, positional := range args.Positional {
		//parameter := t.params.PositionalRequired[i]
		switch positional.Parameter.Type {
		case task.String:
			argStrings[i] = fmt.Sprintf(`"%v"`, positional.Value)
		case task.Bool, task.Int:
			argStrings[i] = fmt.Sprint(positional.Value)
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
