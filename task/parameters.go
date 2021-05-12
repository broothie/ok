package task

import (
	"fmt"
	"strings"
)

type Parameter struct {
	Name      string
	IsKeyword bool
	Default   interface{}
}

func (p Parameter) IsRequired() bool {
	return p.Default == nil
}

type ParamList []Parameter

type Parameters struct {
	ParamList
	Forward bool

	// PositionalRequired []Parameter
	// PositionalOptional []Parameter
	// KeywordRequired    []Parameter
	// KeywordOptional    []Parameter
}

func (l ParamList) Lookup(name string) (Parameter, bool) {
	for _, param := range l {
		if param.Name == name {
			return param, true
		}
	}

	return Parameter{}, false
}

func (l ParamList) Select(isKeyword, isRequired bool) []Parameter {
	params := make([]Parameter, 0, len(l))
	for _, param := range l {
		if param.IsKeyword == isKeyword && param.IsRequired() == isRequired {
			params = append(params, param)
		}
	}

	return params
}

func (p Parameters) PositionalAt(index int) (Parameter, bool) {
	concatenated := append(p.PositionalRequired, p.PositionalOptional...)
	if index < len(concatenated) {
		return concatenated[index], true
	}

	return Parameter{}, false
}

func (p Parameters) KeywordAt(name string) (Parameter, bool) {
	concatenated := append(p.KeywordRequired, p.KeywordOptional...)
	for _, param := range concatenated {
		if param.Name == name {
			return param, true
		}
	}

	// Try with single character
	if len(name) == 1 {
		for _, param := range concatenated {
			if strings.HasPrefix(param.Name, name[:1]) {
				return param, true
			}
		}
	}

	return Parameter{}, false
}

func (p Parameter) String() string {
	if p.Default == nil {
		return p.Name
	}

	return fmt.Sprintf("%s=%s", p.Name, p.Default)
}

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
