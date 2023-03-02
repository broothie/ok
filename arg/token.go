package arg

import "strings"

type Token string

func (t Token) String() string {
	return string(t)
}

func (t Token) IsLongFlag() bool {
	return strings.HasPrefix(t.String(), "--")
}

func (t Token) IsShortFlag() bool {
	return strings.HasPrefix(t.String(), "-") && !t.IsLongFlag()
}

func (t Token) IsFlag() bool {
	return t.IsLongFlag() || t.IsShortFlag()
}

func (t Token) Dashless() string {
	if t.IsLongFlag() {
		return strings.TrimPrefix(t.String(), "--")
	} else {
		return strings.TrimPrefix(t.String(), "-")
	}
}
