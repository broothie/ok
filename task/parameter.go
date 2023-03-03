package task

import "fmt"

type Type string

const (
	TypeString Type = "string"
	TypeBool   Type = "bool"
	TypeInt    Type = "int"
	TypeFloat  Type = "float"
)

type Parameter struct {
	Name    string
	Type    Type
	Default *string
}

func NewRequired(name string, t Type) Parameter {
	return Parameter{Name: name, Type: t}
}

func NewOptional(name string, t Type, dfault string) Parameter {
	return Parameter{Name: name, Type: t, Default: &dfault}
}

func (p Parameter) String() string {
	if p.IsRequired() {
		return fmt.Sprintf("<%s>", p.Name)
	} else {
		return fmt.Sprintf("--%s %s", p.Name, *p.Default)
	}
}

func (p Parameter) IsRequired() bool {
	return p.Default == nil
}

func (p Parameter) IsOptional() bool {
	return !p.IsRequired()
}
