package table

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollapseColumns(t *testing.T) {
	t.Run("drops fully empty columns and preserves order", func(t *testing.T) {
		rows := [][]string{
			{"a", "", "c", ""},
			{"", "", "d", ""},
			{"b", "", "", ""},
		}

		assert.Equal(t, [][]string{
			{"a", "c"},
			{"", "d"},
			{"b", ""},
		}, CollapseColumns(rows))
	})

	t.Run("keeps column if any row has a value", func(t *testing.T) {
		rows := [][]string{
			{"", "x", ""},
			{"", "", ""},
		}

		assert.Equal(t, [][]string{
			{"x"},
			{""},
		}, CollapseColumns(rows))
	})

	t.Run("empty input", func(t *testing.T) {
		assert.Equal(t, [][]string{}, CollapseColumns([][]string{}))
	})
}

func TestWrite(t *testing.T) {
	t.Run("single row", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, Write(&buf, [][]string{{"a", "b", "c"}}))

		output := buf.String()
		assert.Contains(t, output, "a")
		assert.Contains(t, output, "b")
		assert.Contains(t, output, "c")
	})

	t.Run("multiple rows with alignment", func(t *testing.T) {
		var buf bytes.Buffer
		rows := [][]string{
			{"NAME", "VALUE"},
			{"short", "1"},
			{"longer-name", "2"},
		}
		require.NoError(t, Write(&buf, rows))

		output := buf.String()
		assert.Contains(t, output, "NAME")
		assert.Contains(t, output, "VALUE")
		assert.Contains(t, output, "short")
		assert.Contains(t, output, "longer-name")
	})

	t.Run("empty rows", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, Write(&buf, nil))
		assert.Empty(t, buf.String())
	})
}
