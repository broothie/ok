package argparser

import (
	"github.com/samber/lo"
)

type ArgParser struct {
	version  string
	tokens   []string
	flags    []*Flag
	index    int
	taskName string
}

func New(version string, tokens []string, flags ...*Flag) *ArgParser {
	return &ArgParser{
		version: version,
		tokens:  tokens,
		flags:   flags,
	}
}

func (p *ArgParser) Flag(name string) (*Flag, bool) {
	return lo.Find(p.flags, func(f *Flag) bool { return f.Name == name })
}
