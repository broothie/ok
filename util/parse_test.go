package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractCommentIfPresent(t *testing.T) {
	t.Run("octothorpe present", func(t *testing.T) {
		assert.Equal(t, ExtractCommentIfPresent("# something useful", "#"), "something useful")
	})

	t.Run("comment mismatch", func(t *testing.T) {
		assert.Equal(t, ExtractCommentIfPresent("// something useful", "#"), "")
	})
}
