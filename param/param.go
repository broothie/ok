package param

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

func (p Params) PositionalAt(index int) (Param, bool) {
	if index < len(p.PositionalRequired) {
		return p.PositionalRequired[index], true
	} else if index-len(p.PositionalRequired) < len(p.PositionalOptional) {
		return p.PositionalOptional[index-len(p.PositionalRequired)], true
	}

	return Param{}, false
}
