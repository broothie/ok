package cli

import "strings"

type token string

func (t token) String() string {
	return string(t)
}

func (t token) isLongFlag() bool {
	return strings.HasPrefix(t.String(), "--")
}

func (t token) isShortFlag() bool {
	return strings.HasPrefix(t.String(), "-") && !t.isLongFlag()
}

func (t token) isFlag() bool {
	return t.isLongFlag() || t.isShortFlag()
}

func (t token) dashless() string {
	if t.isLongFlag() {
		return strings.TrimPrefix(t.String(), "--")
	} else {
		return strings.TrimPrefix(t.String(), "-")
	}
}
