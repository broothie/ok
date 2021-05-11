package task

import (
	"strconv"
)

//go:generate stringer -type=Type
type Type int

const (
	Untyped Type = iota
	Bool
	Int
	Float
	String = Untyped
)

func (t Type) Parse(s string) (interface{}, error) {
	switch t {
	case Bool:
		return strconv.ParseBool(s)
	case Int:
		return strconv.Atoi(s)
	case Float:
		return strconv.ParseFloat(s, 64)
	default:
		return s, nil
	}
}
