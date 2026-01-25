package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
