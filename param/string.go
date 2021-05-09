package param

import (
	"fmt"
	"strings"
)

func (p Params) String() string {
	var sections []string

	wrappedPositionalRequireds := make([]string, len(p.PositionalRequired))
	counter := 0
	for _, param := range p.PositionalRequired {
		wrappedPositionalRequireds[counter] = chevronWrap(param.Name)
		counter++
	}
	sections = append(sections, wrappedPositionalRequireds...)

	wrappedPositionalOptionals := make([]string, len(p.PositionalOptional))
	counter = 0
	for _, param := range p.PositionalOptional {
		wrappedPositionalOptionals[counter] = chevronWrap(fmt.Sprintf("%s=%v", param.Name, param.Default))
		counter++
	}
	sections = append(sections, wrappedPositionalOptionals...)

	prefixedKeywordRequireds := make([]string, len(p.KeywordRequired))
	counter = 0
	for _, param := range p.KeywordRequired {
		prefixedKeywordRequireds[counter] = dashPrefix(param.Name)
		counter++
	}
	sections = append(sections, prefixedKeywordRequireds...)

	prefixedKeywordOptionals := make([]string, len(p.KeywordOptional))
	counter = 0
	for _, param := range p.KeywordOptional {
		prefixedKeywordOptionals[counter] = dashPrefix(fmt.Sprintf("%s=%v", param.Name, param.Default))
		counter++
	}
	sections = append(sections, prefixedKeywordOptionals...)

	return strings.Join(sections, " ")
}

func chevronWrap(s string) string {
	return fmt.Sprintf("<%s>", s)
}

func dashPrefix(s string) string {
	if len(s) == 1 {
		return fmt.Sprintf("-%s", s)
	} else {
		return fmt.Sprintf("--%s", s)
	}
}
