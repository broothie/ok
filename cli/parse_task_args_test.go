package cli

import (
	"strings"
	"testing"

	"github.com/broothie/ok/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseTaskArgs(t *testing.T) {
	type testCase struct {
		name        string
		args        string
		params      task.Parameters
		expected    task.Arguments
		expectedErr string
	}

	testCases := []testCase{
		{
			name:   "required and present",
			args:   "task Andrew",
			params: task.Parameters{task.NewRequired("name", task.TypeString)},
			expected: task.Arguments{
				{
					Parameter: task.NewRequired("name", task.TypeString),
					Value:     "Andrew",
				},
			},
		},
		{
			name: "required and missing",
			args: "task",
			params: task.Parameters{
				task.NewRequired("name", task.TypeString),
			},
			expectedErr: "missing required args (given 0, expected 1)",
		},
		{
			name: "short 1",
			args: "task Andrew",
			params: task.Parameters{
				task.NewRequired("first-name", task.TypeString),
				task.NewRequired("last-name", task.TypeString),
			},
			expectedErr: "missing required args (given 1, expected 2)",
		},
		{
			name: "optional and present",
			args: "task --name Andrew",
			params: task.Parameters{
				task.NewOptional("name", task.TypeString, "Ted"),
			},
			expected: task.Arguments{
				{
					Parameter: task.NewOptional("name", task.TypeString, "Ted"),
					Value:     "Andrew",
				},
			},
		},
		{
			name: "optional and missing",
			args: "task",
			params: task.Parameters{
				task.NewOptional("name", task.TypeString, "Ted"),
			},
		},
		{
			name: "optional and invalid",
			args: "task --name",
			params: task.Parameters{
				task.NewOptional("name", task.TypeString, "Ted"),
			},
			expectedErr: `no value provided for "--name"`,
		},
		{
			name: "mixed and successful",
			args: "task --greeting Yo Andrew",
			params: task.Parameters{
				task.NewRequired("name", task.TypeString),
				task.NewOptional("greeting", task.TypeString, "Hello"),
			},
			expected: task.Arguments{
				{
					Parameter: task.NewOptional("greeting", task.TypeString, "Hello"),
					Value:     "Yo",
				},
				{
					Parameter: task.NewRequired("name", task.TypeString),
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
