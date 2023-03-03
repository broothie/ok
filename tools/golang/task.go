package golang

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

var packageReplacer = regexp.MustCompile(`package \w+`)

type Task struct {
	Tool
	name       string
	parameters task.Parameters
	filename   string
	goCode     *string
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Parameters() task.Parameters {
	return t.parameters
}

func (t Task) Run(ctx context.Context, args task.Arguments) error {
	goCode := t.generatedGoCode(args)
	goFile, err := os.CreateTemp(".", "Okfile.*.go")
	if err != nil {
		return errors.Wrap(err, "failed to create temp file")
	}

	defer func() {
		if err := os.Remove(goFile.Name()); err != nil {
			fmt.Println("failed to remove temp file", err)
		}
	}()

	if _, err := goFile.WriteString(goCode); err != nil {
		return errors.Wrap(err, "failed to write to temp file")
	}

	if err := goFile.Close(); err != nil {
		return errors.Wrap(err, "failed to close temp file")
	}

	if err := util.CommandContext(ctx, t.Config().Executable(), "run", goFile.Name()).Run(); err != nil {
		return errors.Wrap(err, "failed to run go command")
	}

	return nil
}

func (t Task) generatedGoCode(args task.Arguments) string {
	var argStrings []string
	for _, arg := range args.Required() {
		switch arg.Type {
		case task.TypeBool, task.TypeInt, task.TypeFloat:
			argStrings = append(argStrings, arg.Value)
		case task.TypeString:
			argStrings = append(argStrings, fmt.Sprintf("%q", arg.Value))
		}
	}

	goCode := packageReplacer.ReplaceAllString(*t.goCode, "package main")
	funcMain := fmt.Sprintf("func main() { %s(%s) }", t.name, strings.Join(argStrings, ", "))
	return fmt.Sprintf("%s\n\n%s", goCode, funcMain)
}
