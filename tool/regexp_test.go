package tool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllWhitespace(t *testing.T) {
	assert.True(t, AllWhitespace(""))
	assert.True(t, AllWhitespace(" "))
	assert.True(t, AllWhitespace("\t"))
	assert.True(t, AllWhitespace("	"))
	assert.True(t, AllWhitespace("\n"))
	assert.True(t, AllWhitespace(`
`))
	assert.True(t, AllWhitespace(`
	 `))
}
