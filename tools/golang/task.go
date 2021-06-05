package golang

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
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
	comment      string
	fileContents *string
}

func (t Task) Comment() string {
	return t.comment
}

func (t Task) Params() task.Parameters {
	return t.params
}

func (t Task) Invoke(args task.Args) task.RunningTask {
	file, err := ioutil.TempFile("", "Okfile-*.go")
	if err != nil {
		util.Warn(ToolName, "failed to write go tempfile: %v", err)
		return nil
	}

	argStrings := make([]string, len(args.Positional))
	for i, positional := range args.Positional {
		switch positional.Parameter.Type {
		case task.String:
			argStrings[i] = fmt.Sprintf(`"%v"`, positional.Value)
		case task.Bool, task.Int, task.Float:
			argStrings[i] = fmt.Sprint(positional.Value)
		}
	}

	if err := tmpl.Execute(file, templateData{
		Source:   *t.fileContents,
		TaskName: t.Name(),
		Args:     argStrings,
	}); err != nil {
		util.Warn(ToolName, "failed to write go template: %v", err)
		return nil
	}

	defer file.Close()

	process := util.Exec(ToolName, "run", file.Name())
	go func() {
		process.Wait()
		os.Remove(file.Name())
	}()

	return process
}
