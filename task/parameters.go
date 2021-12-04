package task

import (
	"fmt"
	"strings"
)

type Parameters struct {
	ParamList
	Forward bool
}

type Parameter struct {
	Name      string
	Type      Type
	Default   interface{}
	IsKeyword bool
}

func (p Parameter) IsRequired() bool {
	return p.Default == nil
}

func (p Parameter) String() string {
	if p.Type == Bool {
		return fmt.Sprintf("--[no-]-%s", p.Name)
	}

	inner := p.Name
	if !p.IsRequired() {
		inner = fmt.Sprintf("%s=%s", p.Name, p.Default)
	}

	wrapper := chevronWrap
	if p.IsKeyword {
		wrapper = dashPrefix
	}

	return wrapper(inner)
}

type ParamList []Parameter

func (l ParamList) ToParameters(forward bool) Parameters {
	return Parameters{ParamList: l, Forward: forward}
}

func (l ParamList) Positional() ParamList {
	positionalParams := make([]Parameter, 0, len(l))
	for _, param := range l {
		if !param.IsKeyword {
			positionalParams = append(positionalParams, param)
		}
	}

	return positionalParams
}

func (l ParamList) Keyword() ParamList {
	positionalParams := make([]Parameter, 0, len(l))
	for _, param := range l {
		if param.IsKeyword {
			positionalParams = append(positionalParams, param)
		}
	}

	return positionalParams
}

func (l ParamList) Required() ParamList {
	requiredParams := make([]Parameter, 0, len(l))
	for _, param := range l {
		if param.IsRequired() {
			requiredParams = append(requiredParams, param)
		}
	}

	return requiredParams
}

func (l ParamList) Optional() ParamList {
	optionalParams := make([]Parameter, 0, len(l))
	for _, param := range l {
		if !param.IsRequired() {
			optionalParams = append(optionalParams, param)
		}
	}

	return optionalParams
}

func (l ParamList) PositionalAt(index int) (Parameter, bool) {
	positional := l.Positional()
	if index < len(positional) {
		return positional[index], true
	}

	return Parameter{}, false
}

func (l ParamList) KeywordAt(name string) (Parameter, bool) {
	keywordParams := l.Keyword()
	for _, param := range keywordParams {
		if param.Name == name {
			return param, true
		}
	}

	// Try with single character
	if len(name) == 1 {
		for _, param := range keywordParams {
			if strings.HasPrefix(param.Name, name[:1]) {
				return param, true
			}
		}
	}

	return Parameter{}, false
}

func (l ParamList) String() string {
	return strings.Trim(fmt.Sprint([]Parameter(l)), "[]")
}

func chevronWrap(v interface{}) string {
	return fmt.Sprintf("<%s>", v)
}

func dashPrefix(v interface{}) string {
	return fmt.Sprintf("--%s", v)
}
