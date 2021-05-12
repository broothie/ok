package tool

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllWhitespace(t *testing.T) {
	assert.True(t, AllWhitespace(""))
	assert.True(t, AllWhitespace(" "))
	assert.True(t, AllWhitespace("\t"))
	assert.True(t, AllWhitespace("	"))
	assert.True(t, AllWhitespace("\n  "))
	assert.True(t, AllWhitespace(" \n \t "))
}

func TestSplitOnWhitespace(t *testing.T) {
	assert.Equal(t, []string{""}, SplitOnWhitespace(""))
	assert.Equal(t, []string{"something", "here"}, SplitOnWhitespace("something here"))
	assert.Equal(t, []string{"something", "here"}, SplitOnWhitespace("something\nhere"))
	assert.Equal(t, []string{"something", "here"}, SplitOnWhitespace("something  here"))
	assert.Equal(t, []string{"something", "here"}, SplitOnWhitespace("something  \nhere"))
	assert.Equal(t, []string{"something", "here"}, SplitOnWhitespace("something\there"))
	assert.Equal(t, []string{"something", "here"}, SplitOnWhitespace("something\t here"))
	assert.Equal(t, []string{"something", "here"}, SplitOnWhitespace("something\t\nhere"))
}

func TestSplitOnCommas(t *testing.T) {
	assert.Equal(t, []string{""}, SplitOnCommas(""))
	assert.Equal(t, []string{"", ""}, SplitOnCommas(","))
	assert.Equal(t, []string{"", ""}, SplitOnCommas(", "))
	assert.Equal(t, []string{"", ""}, SplitOnCommas(" , "))
	assert.Equal(t, []string{"", ""}, SplitOnCommas("\n, "))
	assert.Equal(t, []string{"something", "here"}, SplitOnCommas("something,here"))
	assert.Equal(t, []string{"something", "here"}, SplitOnCommas("something, here"))
	assert.Equal(t, []string{"something", "here"}, SplitOnCommas("something , here"))
	assert.Equal(t, []string{"something", "here"}, SplitOnCommas("something\n, here"))
	assert.Equal(t, []string{"something", "here"}, SplitOnCommas("something,\there"))
}

func TestNamedRegexpResult(t *testing.T) {
	re := regexp.MustCompile(`(?P<name>\w+)\s*\((?P<params>.*?)\)`)
	assert.Nil(t, NamedRegexpResult("", re))
	assert.Equal(t, map[string]string{"name": "build", "params": ""}, NamedRegexpResult("build()", re))
	assert.Equal(t, map[string]string{"name": "build", "params": "something"}, NamedRegexpResult("build(something)", re))
}
