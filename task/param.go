package task

import (
	"fmt"
	"strings"
)

type Parameter struct {
	Name    string
	Default interface{}
	Type    Type
}

type Parameters struct {
	PositionalRequired []Parameter
	PositionalOptional []Parameter
	KeywordRequired    []Parameter
	KeywordOptional    []Parameter
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
