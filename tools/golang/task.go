package golang

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
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

func (t Task) Invoke(args task.Args) (task.RunningTask, error) {
	file, err := ioutil.TempFile("", "Okfile-*.go")
	if err != nil {
		return nil, errors.Wrap(err, "failed to write go tempfile")
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
		return nil, errors.Wrap(err, "failed to write go template")
	}

	defer file.Close()

	process, err := util.Exec(ToolName, "run", file.Name())
	if err != nil {
		return nil, err
	}

	go func() {
		process.Wait()
		os.Remove(file.Name())
	}()

	return process, nil
}
