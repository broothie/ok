package task

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Type string

const (
	TypeString Type = "string"
	TypeBool   Type = "bool"
	TypeInt    Type = "int"
	TypeFloat  Type = "float"
)

type Parameters []Parameter

func (p Parameters) Get(index int) (Parameter, bool) {
	if index >= len(p) {
		return Parameter{}, false
	}

	return p[index], true
}

func (p Parameters) Filter(predicate func(Parameter, int) bool) Parameters {
	return lo.Filter(p, predicate)
}

func (p Parameters) Required() Parameters {
	return p.Filter(func(param Parameter, _ int) bool { return param.IsRequired() })
}

func (p Parameters) Optional() Parameters {
	return p.Filter(func(param Parameter, _ int) bool { return param.IsOptional() })
}

func (p Parameters) Find(predicate func(Parameter) bool) (Parameter, bool) {
	return lo.Find(p, predicate)
}

func (p Parameters) String() string {
	var fields []string
	for _, param := range p {
		if param.IsRequired() {
			fields = append(fields, fmt.Sprintf("<%s>", param.Name))
		} else {
			fields = append(fields, fmt.Sprintf("--%s=%s", param.Name, *param.Default))
		}
	}

	return strings.Join(fields, " ")
}
