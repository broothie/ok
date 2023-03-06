package util

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintTable(t *testing.T) {
	table := [][]string{
		{"KEY", "VALUE"},
		{"a", "b"},
	}

	var buf bytes.Buffer
	assert.NoError(t, PrintTable(&buf, table, 3))

	assert.Equal(t,
		"KEY   VALUE\n"+
			"a     b\n",
		buf.String(),
	)
}
