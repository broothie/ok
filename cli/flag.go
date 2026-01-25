package cli

import (
	"strconv"
	"time"

	"github.com/bobg/errors"
)

type FlagType string

const (
	FlagTypeBool     = "bool"
	FlagTypeString   = "string"
	FlagTypeDuration = "duration"
)

type Flag struct {
	Name         string
	Type         FlagType
	Help         string
	Aliases      []string
	Shorts       []rune
	DefaultValue any

	value any
}

func (f *Flag) Value() any {
	if f.value != nil {
		return f.value
	}

	return f.DefaultValue
}

func (f *Flag) IsBool() bool {
	return f.Type == FlagTypeBool
}

func (f *Flag) Parse(s string) error {
	var err error
	switch f.Type {
	case FlagTypeBool:
		f.value, err = strconv.ParseBool(s)

	case FlagTypeString:
		f.value = s

	case FlagTypeDuration:
		f.value, err = time.ParseDuration(s)

	default:
		return errors.Errorf("invalid flag type %q", f.Type)
	}

	return err
}
