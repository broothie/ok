package cli

import (
	"testing"

	"github.com/broothie/ok/ok"
	"github.com/broothie/ok/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseArgs(t *testing.T) {
	taskName := "asdf"

	t.Run("happy path", func(t *testing.T) {
		apple := "mac"
		banana := "green"
		cherry := "red"
		durian := "stinky"
		parser, _ := parserWithArgsAndOptionsParsed(t, taskName, apple, banana, "-d", durian, "--cherry", cherry)

		appleParam := task.Parameter{Name: "apple"}
		bananaParam := task.Parameter{Name: "banana", Default: "yellow"}
		cherryParam := task.Parameter{Name: "cherry", IsKeyword: true}
		durianParam := task.Parameter{Name: "durian", Default: "smelly", IsKeyword: true}

		expectedArgs := task.Args{
			Positional: []task.Arg{
				{Parameter: appleParam, Value: apple},
				{Parameter: bananaParam, Value: banana},
			},
			Keyword: map[string]task.Arg{
				"durian": {Parameter: durianParam, Value: durian},
				"cherry": {Parameter: cherryParam, Value: cherry},
			},
		}

		params := task.ParamList{appleParam, bananaParam, cherryParam, durianParam}.ToParameters(false)
		actualArgs, err := parser.ParseArgs(params)

		assert.NoError(t, err)
		assert.Equal(t, expectedArgs, actualArgs)
	})

	t.Run("positional mismatch", func(t *testing.T) {
		t.Run("too few", func(t *testing.T) {
			parser, _ := parserWithArgsAndOptionsParsed(t, taskName, "1")
			_, err := parser.ParseArgs(task.ParamList{{Name: "a"}, {Name: "b"}}.ToParameters(false))
			assert.EqualError(t, err, "missing positional args: [b]")
		})

		t.Run("too many", func(t *testing.T) {
			t.Run("without optional params", func(t *testing.T) {
				parser, _ := parserWithArgsAndOptionsParsed(t, taskName, "1", "2")
				_, err := parser.ParseArgs(task.ParamList{{Name: "a"}}.ToParameters(false))
				assert.EqualError(t, err, "too many positional args provided, expected max of 1")
			})

			t.Run("with optional params", func(t *testing.T) {
				parser, _ := parserWithArgsAndOptionsParsed(t, taskName, "1", "2", "3")
				_, err := parser.ParseArgs(task.ParamList{{Name: "a"}, {Name: "b", Default: "2"}}.ToParameters(false))
				assert.EqualError(t, err, "too many positional args provided, expected max of 2")
			})

			t.Run("with only optional params", func(t *testing.T) {
				parser, _ := parserWithArgsAndOptionsParsed(t, taskName, "1", "2", "3")
				_, err := parser.ParseArgs(task.ParamList{{Name: "a", Default: "1"}}.ToParameters(false))
				assert.EqualError(t, err, "too many positional args provided, expected max of 1")
			})
		})
	})
}

func parserWithArgsAndOptionsParsed(t *testing.T, args ...string) (*Parser, ok.Options) {
	parser := parserWithArgs(args...)
	options, err := parser.ParseOptions()
	require.NoError(t, err)

	return parser, options
}
