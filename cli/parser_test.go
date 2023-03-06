package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestParser() *parser {
	return newParser([]string{"-f", "file", "task", "arg", "--key", "value"})
}

func TestParser_advance(t *testing.T) {
	parser := newTestParser()
	assert.Equal(t, 0, parser.index)

	parser.advance(2)
	assert.Equal(t, 2, parser.index)
}

func TestParser_token(t *testing.T) {
	parser := newTestParser()
	token, _ := parser.token(0)
	assert.Equal(t, "-f", token.String())

	token, _ = parser.token(3)
	assert.Equal(t, "arg", token.String())
}

func TestParser_peek(t *testing.T) {
	parser := newTestParser()
	peek, _ := parser.peek(3)
	assert.Equal(t, "arg", peek.String())
}

func TestParser_current(t *testing.T) {
	parser := newTestParser()
	current, _ := parser.current()
	assert.Equal(t, "-f", current.String())
}

func TestParser_next(t *testing.T) {
	parser := newTestParser()
	next, _ := parser.next()
	assert.Equal(t, "file", next.String())
}

func TestParser_isExhausted(t *testing.T) {
	parser := newTestParser()
	assert.False(t, parser.isExhausted())

	parser.advance(5)
	assert.False(t, parser.isExhausted())

	parser.advance(1)
	assert.True(t, parser.isExhausted())

	parser.advance(1)
	assert.True(t, parser.isExhausted())
}
