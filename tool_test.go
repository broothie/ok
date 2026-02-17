package ok

import (
	"bytes"
	"errors"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bmatcuk/doublestar/v4"
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

func TestSkipDirsFS(t *testing.T) {
	inner := fstest.MapFS{
		"scripts/deploy.sh":                    &fstest.MapFile{},
		"node_modules/foo/install.sh":          &fstest.MapFile{},
		"node_modules/bar/postinstall.sh":      &fstest.MapFile{},
		".git/hooks/pre-commit":                &fstest.MapFile{},
		"sub/node_modules/pkg/run.sh":          &fstest.MapFile{},
		"vendor/scripts/build.sh":              &fstest.MapFile{},
		"top.sh":                               &fstest.MapFile{},
	}

	fsys := skipDirsFS{inner: inner}

	matches, err := doublestar.Glob(fsys, "**/*.sh")
	require.NoError(t, err)

	assert.Contains(t, matches, "scripts/deploy.sh")
	assert.Contains(t, matches, "top.sh")
	assert.Contains(t, matches, "vendor/scripts/build.sh")

	for _, match := range matches {
		assert.NotContains(t, match, "node_modules")
		assert.NotContains(t, match, ".git")
	}
}

func TestSkipDirsFS_ReadDir(t *testing.T) {
	inner := fstest.MapFS{
		"node_modules/foo/bar.sh": &fstest.MapFile{},
		".git/hooks/pre-commit":   &fstest.MapFile{},
		".hg/store/data":          &fstest.MapFile{},
		".svn/entries":            &fstest.MapFile{},
		"scripts/deploy.sh":       &fstest.MapFile{},
		"README.md":               &fstest.MapFile{},
	}

	fsys := skipDirsFS{inner: inner}

	entries, err := fsys.ReadDir(".")
	require.NoError(t, err)

	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}

	assert.Contains(t, names, "scripts")
	assert.Contains(t, names, "README.md")
	assert.NotContains(t, names, "node_modules")
	assert.NotContains(t, names, ".git")
	assert.NotContains(t, names, ".hg")
	assert.NotContains(t, names, ".svn")
}

func TestSkipDirsFS_Open(t *testing.T) {
	inner := fstest.MapFS{
		"scripts/deploy.sh": &fstest.MapFile{Data: []byte("#!/bin/bash")},
	}

	fsys := skipDirsFS{inner: inner}

	f, err := fsys.Open("scripts/deploy.sh")
	require.NoError(t, err)
	require.NoError(t, f.Close())

	_, err = fsys.Open("nonexistent")
	assert.ErrorIs(t, err, fs.ErrNotExist)
}
