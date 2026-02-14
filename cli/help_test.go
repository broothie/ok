package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_WriteHelp(t *testing.T) {
	t.Run("includes version", func(t *testing.T) {
		p := NewParser("v1.2.3", nil,
			&Flag{Name: "verbose", Type: FlagTypeBool, Shorts: []rune{'v'}, DefaultValue: false, Help: "Enable verbose output."},
		)

		var buf bytes.Buffer
		require.NoError(t, p.WriteHelp(&buf))

		output := buf.String()
		assert.Contains(t, output, "v1.2.3")
		assert.Contains(t, output, "Usage:")
		assert.Contains(t, output, "--verbose")
		assert.Contains(t, output, "-v")
		assert.Contains(t, output, "Enable verbose output.")
	})

	t.Run("includes aliases", func(t *testing.T) {
		p := NewParser("v0.1.0", nil,
			&Flag{Name: "filter-tools", Type: FlagTypeString, Aliases: []string{"ft"}, DefaultValue: "", Help: "Filter tools."},
		)

		var buf bytes.Buffer
		require.NoError(t, p.WriteHelp(&buf))

		output := buf.String()
		assert.Contains(t, output, "--filter-tools")
		assert.Contains(t, output, "--ft")
	})

	t.Run("no flags", func(t *testing.T) {
		p := NewParser("v0.0.1", nil)

		var buf bytes.Buffer
		require.NoError(t, p.WriteHelp(&buf))

		output := buf.String()
		assert.Contains(t, output, "v0.0.1")
		assert.Contains(t, output, "Usage:")
	})
}
