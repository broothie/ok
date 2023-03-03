package cli

import "sync"

type CLI struct {
	parser *optionParser
	flags  []flag

	optionsOnce *sync.Once
	options     Options
}

func New(args []string) *CLI {
	return &CLI{
		parser:      newOptionParser(args),
		flags:       flags(),
		optionsOnce: new(sync.Once),
		options:     Options{},
	}
}
