package task

import (
	"fmt"
	"strings"
)

func (p Parameters) String() string {
	var sections []string

	wrappedPositionalRequireds := make([]string, len(p.PositionalRequired))
	counter := 0
	for _, param := range p.PositionalRequired {
		wrappedPositionalRequireds[counter] = chevronWrap(param)
		counter++
	}
	sections = append(sections, wrappedPositionalRequireds...)

	wrappedPositionalOptionals := make([]string, len(p.PositionalOptional))
	counter = 0
	for _, param := range p.PositionalOptional {
		wrappedPositionalOptionals[counter] = chevronWrap(param)
		counter++
	}
	sections = append(sections, wrappedPositionalOptionals...)

	prefixedKeywordRequireds := make([]string, len(p.KeywordRequired))
	counter = 0
	for _, param := range p.KeywordRequired {
		prefixedKeywordRequireds[counter] = dashPrefix(param)
		counter++
	}
	sections = append(sections, prefixedKeywordRequireds...)

	prefixedKeywordOptionals := make([]string, len(p.KeywordOptional))
	counter = 0
	for _, param := range p.KeywordOptional {
		prefixedKeywordOptionals[counter] = dashPrefix(param)
		counter++
	}
	sections = append(sections, prefixedKeywordOptionals...)

	return strings.Join(sections, " ")
}

func chevronWrap(v interface{}) string {
	return fmt.Sprintf("<%s>", v)
}

func dashPrefix(v interface{}) string {
	return fmt.Sprintf("--%s", v)
}
