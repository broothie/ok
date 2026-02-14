package ok

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	v := Version()
	assert.NotEmpty(t, v)
	assert.NotContains(t, v, "\n", "version should not contain newlines")
	assert.Equal(t, v, Version(), "version should be consistent across calls")
}
