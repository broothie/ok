package parameter

type Parameter struct {
	Name    string
	Type    Type
	Default *string
}

func (p Parameter) IsRequired() bool {
	return p.Default == nil
}

func (p Parameter) IsOptional() bool {
	return !p.IsRequired()
}
