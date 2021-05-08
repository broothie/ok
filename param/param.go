package param

type Type int

const (
	Untyped Type = iota
	Bool
	Int
	String
	ListUntyped
	ListBool
	ListInt
	ListString
)

type Param struct {
	Name    string
	Type    Type
	Default interface{}
}

type Params struct {
	PositionalRequired []Param
	PositionalOptional []Param
	KeywordRequired    []Param
	KeywordOptional    []Param
}
