package ok

import (
	"testing"

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
		cherryParam := task.Parameter{Name: "cherry"}
		durianParam := task.Parameter{Name: "durian", Default: "smelly"}

		expectedArgs := task.Args{
			Positional: []task.Arg{{Parameter: appleParam, Value: apple}, {Parameter: bananaParam, Value: banana}},
			Keyword: map[string]task.Arg{
				"durian": {Parameter: durianParam, Value: durian},
				"cherry": {Parameter: cherryParam, Value: cherry},
			},
		}

		actualArgs, err := parser.ParseArgs(task.Parameters{
			PositionalRequired: []task.Parameter{appleParam},
			PositionalOptional: []task.Parameter{bananaParam},
			KeywordRequired:    []task.Parameter{cherryParam},
			KeywordOptional:    []task.Parameter{durianParam},
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedArgs, actualArgs)
	})

	t.Run("positional mismatch", func(t *testing.T) {
		t.Run("too few", func(t *testing.T) {
			parser, _ := parserWithArgsAndOptionsParsed(t, taskName, "1")

			_, err := parser.ParseArgs(task.Parameters{
				PositionalRequired: []task.Parameter{{Name: "a"}, {Name: "b"}},
			})

			assert.EqualError(t, err, "missing positional args: [b]")
		})

		t.Run("too many", func(t *testing.T) {
			t.Run("without optional params", func(t *testing.T) {
				parser, _ := parserWithArgsAndOptionsParsed(t, taskName, "1", "2")

				_, err := parser.ParseArgs(task.Parameters{
					PositionalRequired: []task.Parameter{{Name: "a"}},
				})

				assert.EqualError(t, err, "too many positional args provided, expected max of 1")
			})

			t.Run("with optional params", func(t *testing.T) {
				parser, _ := parserWithArgsAndOptionsParsed(t, taskName, "1", "2", "3")

				_, err := parser.ParseArgs(task.Parameters{
					PositionalRequired: []task.Parameter{{Name: "a"}},
					PositionalOptional: []task.Parameter{{Name: "b"}},
				})

				assert.EqualError(t, err, "too many positional args provided, expected max of 2")
			})

			t.Run("with only optional params", func(t *testing.T) {
				parser, _ := parserWithArgsAndOptionsParsed(t, taskName, "1", "2", "3")

				_, err := parser.ParseArgs(task.Parameters{
					PositionalOptional: []task.Parameter{{Name: "a"}},
				})

				assert.EqualError(t, err, "too many positional args provided, expected max of 1")
			})
		})
	})
}

func parserWithArgsAndOptionsParsed(t *testing.T, args ...string) (*Parser, Options) {
	parser := NewParser(args)
	options, err := parser.ParseOptions()
	require.NoError(t, err)

	return parser, options
}
