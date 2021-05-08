package param

import (
	"fmt"
	"strings"
)

func (a Params) String() string {
	var sections []string

	wrappedPositionalRequireds := make([]string, len(a.PositionalRequired))
	counter := 0
	for _, param := range a.PositionalRequired {
		wrappedPositionalRequireds[counter] = chevronWrap(param.Name)
		counter++
	}
	sections = append(sections, wrappedPositionalRequireds...)

	wrappedPositionalOptionals := make([]string, len(a.PositionalOptional))
	counter = 0
	for _, param := range a.PositionalOptional {
		wrappedPositionalOptionals[counter] = chevronWrap(fmt.Sprintf("%s=%v", param.Name, param.Default))
		counter++
	}
	sections = append(sections, wrappedPositionalOptionals...)

	prefixedKeywordRequireds := make([]string, len(a.KeywordRequired))
	counter = 0
	for _, param := range a.KeywordRequired {
		prefixedKeywordRequireds[counter] = dashPrefix(param.Name)
		counter++
	}
	sections = append(sections, prefixedKeywordRequireds...)

	prefixedKeywordOptionals := make([]string, len(a.KeywordOptional))
	counter = 0
	for _, param := range a.KeywordOptional {
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
