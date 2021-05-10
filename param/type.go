package param

import (
	"fmt"
	"strconv"
)

//go:generate stringer -type=Type
type Type int

const (
	Untyped Type = iota
	Bool
	Int
	Float
	String
)

func (t Type) CastString(s string) (interface{}, error) {
	switch t {
	case Untyped:
		return s, nil
	case Bool:
		return strconv.ParseBool(s)
	case Int:
		return strconv.Atoi(s)
	case Float:
		return strconv.ParseFloat(s, 64)
	case String:
		return s, nil
	default:
		return "", fmt.Errorf("invalid param type '%v'", t)
	}
}
