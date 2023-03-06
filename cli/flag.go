package cli

type applyFunc func(*parser, *Options) error

type flag struct {
	long        string
	short       rune
	valueName   string
	description string
	apply       applyFunc
}

func (f flag) isMatch(arg token) bool {
	if !arg.isFlag() {
		return false
	} else if arg.isLongFlag() {
		return f.long == arg.dashless()
	} else {
		return string(f.short) == arg.dashless()
	}
}

func (f flag) hasShort() bool {
	return f.short != 0
}
