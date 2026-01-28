package ok

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOk_ListTools(t *testing.T) {
	type testCase struct {
		tools            []toolInfo
		expectedContains []string
		expectedOrder    []string
	}

	tests := map[string]testCase{
		"no tools": {
			tools:            nil,
			expectedContains: []string{"TOOL", "EXECUTABLE", "FILES"},
			expectedOrder:    nil,
		},
		"single tool with valid command": {
			tools: []toolInfo{
				{
					Tool:        Tool{Name: "Make"},
					commandPath: "/usr/bin/make",
					filePaths:   []string{"Makefile"},
				},
			},
			expectedContains: []string{"Make", "✔ /usr/bin/make", "Makefile"},
			expectedOrder:    []string{"Make"},
		},
		"single tool with command error": {
			tools: []toolInfo{
				{
					Tool:           Tool{Name: "Make"},
					commandPathErr: errors.New("not found"),
					filePaths:      []string{"Makefile"},
				},
			},
			expectedContains: []string{"Make", "✘ not found", "Makefile"},
			expectedOrder:    []string{"Make"},
		},
		"single tool with file error": {
			tools: []toolInfo{
				{
					Tool:         Tool{Name: "Make"},
					commandPath:  "/usr/bin/make",
					filePathsErr: errors.New("glob error"),
				},
			},
			expectedContains: []string{"Make", "✔ /usr/bin/make", "✘ glob error"},
			expectedOrder:    []string{"Make"},
		},
		"single tool with multiple files": {
			tools: []toolInfo{
				{
					Tool:        Tool{Name: "Make"},
					commandPath: "/usr/bin/make",
					filePaths:   []string{"Makefile", "sub/Makefile"},
				},
			},
			expectedContains: []string{"Make", "✔ /usr/bin/make", "Makefile,sub/Makefile"},
			expectedOrder:    []string{"Make"},
		},
		"multiple tools sorted alphabetically": {
			tools: []toolInfo{
				{
					Tool:        Tool{Name: "Rake"},
					commandPath: "/usr/bin/rake",
					filePaths:   []string{"Rakefile"},
				},
				{
					Tool:        Tool{Name: "Make"},
					commandPath: "/usr/bin/make",
					filePaths:   []string{"Makefile"},
				},
			},
			expectedContains: []string{"Make", "Rake", "Makefile", "Rakefile"},
			expectedOrder:    []string{"Make", "Rake"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			o := &Ok{tools: tc.tools}

			var buf bytes.Buffer
			err := o.ListTools(&buf)

			require.NoError(t, err)

			output := buf.String()
			for _, s := range tc.expectedContains {
				assert.Contains(t, output, s)
			}

			// Verify sort order
			if len(tc.expectedOrder) > 1 {
				for i := 0; i < len(tc.expectedOrder)-1; i++ {
					pos1 := strings.Index(output, tc.expectedOrder[i])
					pos2 := strings.Index(output, tc.expectedOrder[i+1])
					assert.Less(t, pos1, pos2, "expected %q before %q", tc.expectedOrder[i], tc.expectedOrder[i+1])
				}
			}
		})
	}
}
