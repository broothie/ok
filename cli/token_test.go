package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_token(t *testing.T) {
	assert.True(t, token("--flag").isLongFlag())
	assert.False(t, token("-flag").isLongFlag())
	assert.False(t, token("flag").isLongFlag())

	assert.False(t, token("--flag").isShortFlag())
	assert.True(t, token("-flag").isShortFlag())
	assert.False(t, token("flag").isShortFlag())

	assert.True(t, token("--flag").isFlag())
	assert.True(t, token("-flag").isFlag())
	assert.False(t, token("flag").isFlag())

	assert.Equal(t, "flag", token("--flag").dashless())
	assert.Equal(t, "flag", token("-flag").dashless())
	assert.Equal(t, "flag", token("flag").dashless())
}
