package cli

import (
	"strings"
	"testing"

	"github.com/broothie/ok/argument"
	"github.com/broothie/ok/parameter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseTaskArgs(t *testing.T) {
	type testCase struct {
		name        string
		args        string
		params      parameter.Parameters
		expected    argument.Arguments
		expectedErr string
	}

	testCases := []testCase{
		{
			name:   "required and present",
			args:   "task Andrew",
			params: parameter.Parameters{parameter.NewRequired("name", parameter.TypeString)},
			expected: argument.Arguments{
				{
					Parameter: parameter.NewRequired("name", parameter.TypeString),
					Value:     "Andrew",
				},
			},
		},
		{
			name: "required and missing",
			args: "task",
			params: parameter.Parameters{
				parameter.NewRequired("name", parameter.TypeString),
			},
			expectedErr: "missing required args (given 0, expected 1)",
		},
		{
			name: "short 1",
			args: "task Andrew",
			params: parameter.Parameters{
				parameter.NewRequired("first-name", parameter.TypeString),
				parameter.NewRequired("last-name", parameter.TypeString),
			},
			expectedErr: "missing required args (given 1, expected 2)",
		},
		{
			name: "optional and present",
			args: "task --name Andrew",
			params: parameter.Parameters{
				parameter.NewOptional("name", parameter.TypeString, "Ted"),
			},
			expected: argument.Arguments{
				{
					Parameter: parameter.NewOptional("name", parameter.TypeString, "Ted"),
					Value:     "Andrew",
				},
			},
		},
		{
			name: "optional and missing",
			args: "task",
			params: parameter.Parameters{
				parameter.NewOptional("name", parameter.TypeString, "Ted"),
			},
		},
		{
			name: "optional and invalid",
			args: "task --name",
			params: parameter.Parameters{
				parameter.NewOptional("name", parameter.TypeString, "Ted"),
			},
			expectedErr: `no value provided for "--name"`,
		},
		{
			name: "mixed and successful",
			args: "task --greeting Yo Andrew",
			params: parameter.Parameters{
				parameter.NewRequired("name", parameter.TypeString),
				parameter.NewOptional("greeting", parameter.TypeString, "Hello"),
			},
			expected: argument.Arguments{
				{
					Parameter: parameter.NewOptional("greeting", parameter.TypeString, "Hello"),
					Value:     "Yo",
				},
				{
					Parameter: parameter.NewRequired("name", parameter.TypeString),
					Value:     "Andrew",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := New(strings.Fields(testCase.args))
			_, err := parser.ParseOkArgs()
			require.NoError(t, err)

			actual, err := parser.ParseTaskArgs(testCase.params)
			if err != nil {
				assert.EqualError(t, err, testCase.expectedErr)
				return
			}

			assert.Equal(t, testCase.expected, actual)
		})
	}
}
