package ruby

import (
	"fmt"
	"os"
	"strings"

	"github.com/broothie/now/arg"
	"github.com/broothie/now/param"
	"github.com/broothie/now/task"
	"github.com/broothie/now/tool"
)

type Task struct {
	task.Base
	params param.Params
}

func (t Task) Params() param.Params {
	return t.params
}

func (t Task) Invoke(args arg.Task) *os.Process {
	positionalStrings := make([]string, len(args.Positional))
	for i, positional := range args.Positional {
		positionalStrings[i] = fmt.Sprintf(`"%s"`, positional)
	}

	//builder := new(strings.Builder)
	//builder.WriteString(strings.Join())

	return tool.Exec(ToolName,
		"-r", fmt.Sprintf("./%s", t.Filename()),
		"-e", fmt.Sprintf("%s(%s)", t.Name(), strings.Join(positionalStrings, ", ")),
	).Process
}
