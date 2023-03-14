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

func SplatParameters(t Type) Parameters {
	return Parameters{SplatParameter(t)}
}

func SplatParameter(t Type) Parameter {
	return PositionalParameter("...", t)
}

func PositionalParameter(name string, t Type) Parameter {
	return Parameter{Name: name, Type: t}
}

func KeywordParameter(name string, t Type, dfault string) Parameter {
	return Parameter{Name: name, Type: t, Default: &dfault}
}

func (p Parameter) String() string {
	if p.IsSplat() {
		return p.Name
	} else if p.IsPositional() {
		return fmt.Sprintf("<%s>", p.Name)
	} else {
		return fmt.Sprintf("--%s %s", p.Name, *p.Default)
	}
}

func (p Parameter) IsSplat() bool {
	return p.Name == "..."
}

func (p Parameter) IsPositional() bool {
	return p.Default == nil
}

func (p Parameter) IsKeyword() bool {
	return !p.IsPositional()
}
