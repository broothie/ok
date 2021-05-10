package ruby

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/broothie/now/arg"
	"github.com/broothie/now/param"
	"github.com/broothie/now/task"
	"github.com/broothie/now/toolhelp"
)

type Task struct {
	task.Base
	params param.Params
}

func (t Task) Params() param.Params {
	return t.params
}

func (t Task) Invoke(args arg.Args) *os.Process {
	positionalStrings := make([]string, len(args.Positional))
	for i, arg := range args.Positional {
		positionalStrings[i] = processArg(arg.Value.(string))
	}

	keywordEntries := make([]string, len(args.Keyword))
	counter := 0
	for name, arg := range args.Keyword {
		keywordEntries[counter] = fmt.Sprintf("%s: %s", name, processArg(arg.Value.(string)))
		counter++
	}

	builder := new(strings.Builder)
	builder.WriteString(fmt.Sprintf("args = [%s]; ", strings.Join(positionalStrings, ", ")))
	builder.WriteString(fmt.Sprintf("kwargs = {%s}; ", strings.Join(keywordEntries, ", ")))
	builder.WriteString(fmt.Sprintf("%s(*args, **kwargs)", t.Name()))

	return toolhelp.Exec(ToolName,
		"-r", fmt.Sprintf("./%s", t.Filename()),
		"-e", builder.String(),
	).Process
}

func processArg(arg string) string {
	if _, err := strconv.ParseFloat(arg, 64); err == nil {
		return arg
	} else if _, err := strconv.Atoi(arg); err == nil {
		return arg
	} else if _, err := strconv.ParseBool(arg); err == nil {
		return arg
	} else {
		return fmt.Sprintf(`"%s"`, arg)
	}
}
