package golang

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/util"
	"github.com/pkg/errors"
)

var packageReplacer = regexp.MustCompile(`package \w[a-zA-Z0-9]+`)

type Task struct {
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
	sum := sha256.Sum256([]byte(goCode))
	goFile, err := os.Create(fmt.Sprintf("Okfile.%s.go", hex.EncodeToString(sum[:])))
	if err != nil {
		return errors.Wrap(err, "failed to create temp file")
	}

	defer func() {
		if err := os.Remove(goFile.Name()); err != nil {
			logger.Log.Printf("failed to remove temp file: %v", err)
		}
	}()

	if _, err := goFile.WriteString(goCode); err != nil {
		return errors.Wrap(err, "failed to write to temp file")
	}

	if err := goFile.Close(); err != nil {
		return errors.Wrap(err, "failed to close temp file")
	}

	if err := util.CommandContext(ctx, "go", "run", goFile.Name()).Run(); err != nil {
		return errors.Wrap(err, "failed to run go task")
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
