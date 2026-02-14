package cli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlag_IsBool(t *testing.T) {
	tests := map[string]struct {
		flagType FlagType
		expected bool
	}{
		"bool":     {FlagTypeBool, true},
		"string":   {FlagTypeString, false},
		"duration": {FlagTypeDuration, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &Flag{Type: tc.flagType}
			assert.Equal(t, tc.expected, f.IsBool())
		})
	}
}

func TestFlag_Value(t *testing.T) {
	t.Run("returns default when no value set", func(t *testing.T) {
		f := &Flag{DefaultValue: "default"}
		assert.Equal(t, "default", f.Value())
	})

	t.Run("returns parsed value after Parse", func(t *testing.T) {
		f := &Flag{Type: FlagTypeString, DefaultValue: "default"}
		require.NoError(t, f.Parse("custom"))
		assert.Equal(t, "custom", f.Value())
	})
}

func TestFlag_Parse(t *testing.T) {
	t.Run("bool true values", func(t *testing.T) {
		for _, input := range []string{"true", "TRUE", "True", "1", "t", "T"} {
			f := &Flag{Type: FlagTypeBool}
			require.NoError(t, f.Parse(input))
			assert.Equal(t, true, f.Value(), "input: %s", input)
		}
	})

	t.Run("bool false values", func(t *testing.T) {
		for _, input := range []string{"false", "FALSE", "False", "0", "f", "F"} {
			f := &Flag{Type: FlagTypeBool}
			require.NoError(t, f.Parse(input))
			assert.Equal(t, false, f.Value(), "input: %s", input)
		}
	})

	t.Run("bool invalid", func(t *testing.T) {
		f := &Flag{Type: FlagTypeBool}
		assert.Error(t, f.Parse("notabool"))
	})

	t.Run("string", func(t *testing.T) {
		f := &Flag{Type: FlagTypeString}
		require.NoError(t, f.Parse("hello"))
		assert.Equal(t, "hello", f.Value())
	})

	t.Run("string empty", func(t *testing.T) {
		f := &Flag{Type: FlagTypeString}
		require.NoError(t, f.Parse(""))
		assert.Equal(t, "", f.Value())
	})

	t.Run("duration", func(t *testing.T) {
		f := &Flag{Type: FlagTypeDuration}
		require.NoError(t, f.Parse("5s"))
		assert.Equal(t, 5*time.Second, f.Value())
	})

	t.Run("duration complex", func(t *testing.T) {
		f := &Flag{Type: FlagTypeDuration}
		require.NoError(t, f.Parse("1m30s"))
		assert.Equal(t, 90*time.Second, f.Value())
	})

	t.Run("duration invalid", func(t *testing.T) {
		f := &Flag{Type: FlagTypeDuration}
		assert.Error(t, f.Parse("notaduration"))
	})

	t.Run("invalid flag type", func(t *testing.T) {
		f := &Flag{Type: "invalid"}
		assert.Error(t, f.Parse("anything"))
	})
}
