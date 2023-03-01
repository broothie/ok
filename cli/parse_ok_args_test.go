package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_ParseOkArgs(t *testing.T) {
	type testCase struct {
		name        string
		args        string
		expected    OkArgs
		expectedErr string
	}

	testCases := []testCase{
		{
			name: "-h",
			args: "-h",
			expected: OkArgs{
				Help:      true,
				Version:   false,
				ListTools: false,
				Watches:   nil,
				TaskName:  "",
			},
		},
		{
			name: "--help",
			args: "--help",
			expected: OkArgs{
				Help:      true,
				Version:   false,
				ListTools: false,
				Watches:   nil,
				TaskName:  "",
			},
		},
		{
			name: "-V",
			args: "-V",
			expected: OkArgs{
				Help:      false,
				Version:   true,
				ListTools: false,
				Watches:   nil,
				TaskName:  "",
			},
		},
		{
			name: "--version",
			args: "--version",
			expected: OkArgs{
				Help:      false,
				Version:   true,
				ListTools: false,
				Watches:   nil,
				TaskName:  "",
			},
		},
		{
			name: "--tools",
			args: "--tools",
			expected: OkArgs{
				Help:      false,
				Version:   false,
				ListTools: true,
				Watches:   nil,
				TaskName:  "",
			},
		},
		{
			name: "-w README.md",
			args: "-w README.md",
			expected: OkArgs{
				Help:      false,
				Version:   false,
				ListTools: false,
				Watches:   []string{"README.md"},
				TaskName:  "",
			},
		},
		{
			name: "--watch README.md",
			args: "--watch README.md",
			expected: OkArgs{
				Help:      false,
				Version:   false,
				ListTools: false,
				Watches:   []string{"README.md"},
				TaskName:  "",
			},
		},
		{
			name: "-w VERSION --watch README.md",
			args: "-w VERSION --watch README.md",
			expected: OkArgs{
				Help:      false,
				Version:   false,
				ListTools: false,
				Watches:   []string{"VERSION", "README.md"},
				TaskName:  "",
			},
		},
		{
			name: "task",
			args: "greet",
			expected: OkArgs{
				Help:      false,
				Version:   false,
				ListTools: false,
				Watches:   nil,
				TaskName:  "greet",
			},
		},
		{
			name:        "-x",
			args:        "-x",
			expectedErr: `invalid option "-x"`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			parser := New(strings.Fields(testCase.args))

			actual, err := parser.ParseOkArgs()
			if err != nil {
				assert.EqualError(t, err, testCase.expectedErr)
				return
			}

			assert.Equal(t, testCase.expected, actual)
		})
	}
}
