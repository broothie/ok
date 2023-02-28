package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCommaArgList(t *testing.T) {
	assert.Equal(t, []string{"some", "params"}, SplitCommaParamList("some,params"))
	assert.Equal(t, []string{"some", "params"}, SplitCommaParamList("some, params"))
}
