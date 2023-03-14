package cli

import (
	"testing"

	"github.com/broothie/ok/task"
	"github.com/stretchr/testify/assert"
)

func TestCLI_ParseParameters(t *testing.T) {
	type testCase struct {
		name              string
		args              []string
		parameters        task.Parameters
		expectedArguments task.Arguments
		expectedError     error
	}

	testCases := []testCase{
		{
			name: "positional arg",
			args: []string{"hello"},
			parameters: task.Parameters{
				task.PositionalParameter("greeting", task.TypeString),
			},
			expectedArguments: task.Arguments{
				{
					Parameter: task.PositionalParameter("greeting", task.TypeString),
					Value:     "hello",
				},
			},
		},
		{
			name: "keyword arg",
			args: []string{"--greeting", "hello"},
			parameters: task.Parameters{
				task.KeywordParameter("greeting", task.TypeString, "hi"),
			},
			expectedArguments: task.Arguments{
				{
					Parameter: task.KeywordParameter("greeting", task.TypeString, "hi"),
					Value:     "hello",
				},
			},
		},
		{
			name: "positional, 2 keywords, one argued",
			args: []string{"--greeting", "hello", "andrew"},
			parameters: task.Parameters{
				task.PositionalParameter("name", task.TypeString),
				task.KeywordParameter("excited", task.TypeBool, "false"),
				task.KeywordParameter("greeting", task.TypeString, "hi"),
			},
			expectedArguments: task.Arguments{
				{
					Parameter: task.KeywordParameter("greeting", task.TypeString, "hi"),
					Value:     "hello",
				},
				{
					Parameter: task.PositionalParameter("name", task.TypeString),
					Value:     "andrew",
				},
			},
		},
		{
			name:              "splat, none argued",
			args:              []string{},
			parameters:        task.SplatParameters(task.TypeString),
			expectedArguments: nil,
		},
		{
			name:       "splat, two argued",
			args:       []string{"hi", "there"},
			parameters: task.SplatParameters(task.TypeString),
			expectedArguments: task.Arguments{
				{
					Parameter: task.SplatParameter(task.TypeString),
					Value:     "hi",
				},
				{
					Parameter: task.SplatParameter(task.TypeString),
					Value:     "there",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cli := New(testCase.args)

			actualArguments, err := cli.ParseParameters(testCase.parameters)
			if testCase.expectedError != nil {
				assert.EqualError(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedArguments, actualArguments)
			}
		})
	}
}
