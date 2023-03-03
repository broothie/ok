package task

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

func (p Parameter) IsRequired() bool {
	return p.Default == nil
}

func (p Parameter) IsOptional() bool {
	return !p.IsRequired()
}
