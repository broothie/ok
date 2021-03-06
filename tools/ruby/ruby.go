package ruby

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/task"
	"github.com/broothie/ok/tools/ez"
	"github.com/broothie/ok/util"
)

type Config struct {
	Bundler *bool `toml:"bundler"`
	Rails   *bool `toml:"rails"`
}

var (
	methodMatcher     = regexp.MustCompile(`(?m)^def\s+(?P<taskName>\w+)\(?(?P<params>.*?)\)?$`)
	positionalMatcher = regexp.MustCompile(`^(?P<paramName>\w+)(?:\s*=\s*(?P<default>.*))?$`)
	keywordMatcher    = regexp.MustCompile(`^(?P<paramName>\w+):(?:\s*(?P<default>.*))?$`)

	Ruby = ez.Tool{
		ToolName:             "ruby",
		CommandName:          "ruby",
		ToolFilename:         "Okfile.rb",
		TaskMatcher:          methodMatcher,
		CommentPrefixMatcher: util.OctothorpePrefixMatcher,
		ToolConfig:           new(Config),
		ParamParser: func(paramString string) (task.Parameters, error) {
			return paramListFromParamString(paramString), nil
		},
		Invoke: func(task ez.Task, args task.Args) (task.RunningTask, error) {
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

			methodArgs := append(positionalStrings, keywordEntries...)
			script := fmt.Sprintf("%s; %s(%s)", *task.FileContents, task.Name(), strings.Join(methodArgs, ", "))
			execArgs := []string{"ruby", "-e", script}
			config := task.Tool.ToolConfig.(*Config) // TODO: Type safe way to do this?
			if config.Rails != nil && *config.Rails {
				execArgs = []string{"rails", "runner", script}
			}

			if config.Bundler != nil && *config.Bundler {
				execArgs = append([]string{"bin/bundle", "exec"}, execArgs...)
			}

			return util.Exec(execArgs[0], execArgs[1:]...)
		},
	}
)

func paramListFromParamString(paramsString string) task.Parameters {
	paramStrings := util.SplitOnCommas(paramsString)
	if len(paramStrings) == 1 && util.AllWhitespace(paramStrings[0]) {
		return task.Parameters{}
	}

	var params task.Parameters
	for _, paramString := range paramStrings {
		var re *regexp.Regexp
		var isKeyword bool

		if positionalMatcher.MatchString(paramString) {
			re = positionalMatcher
			isKeyword = false
		} else if keywordMatcher.MatchString(paramString) {
			re = keywordMatcher
			isKeyword = true
		} else {
			logger.Tool("ruby").Printf("error parsing param %q", paramString)
			continue
		}

		result := util.NamedRegexpResult(paramString, re)

		var defaultValue interface{}
		defaultString, defaultPresent := result["default"]
		if defaultPresent && defaultString != "" {
			defaultValue = defaultString
		}

		defaultString = strings.TrimSpace(defaultString)
		if strings.HasPrefix(defaultString, "'") || strings.HasPrefix(defaultString, `"`) {
			defaultValue = strings.Trim(defaultString, `'"`)
		}

		params.ParamList = append(params.ParamList, task.Parameter{
			Name:      result["paramName"],
			Type:      task.Untyped,
			Default:   defaultValue,
			IsKeyword: isKeyword,
		})
	}

	return params
}

func processArg(arg string) string {
	if _, err := strconv.ParseFloat(arg, 64); err == nil {
		return arg
	} else if _, err := strconv.Atoi(arg); err == nil {
		return arg
	} else if _, err := strconv.ParseBool(arg); err == nil {
		return arg
	} else {
		return strconv.Quote(arg)
	}
}
