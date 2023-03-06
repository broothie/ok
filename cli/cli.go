package cli

import "sync"

type CLI struct {
	parser *parser
	flags  []flag

	optionsOnce *sync.Once
	options     Options
}

func New(args []string) *CLI {
	return &CLI{
		parser:      newParser(args),
		flags:       flags(),
		optionsOnce: new(sync.Once),
		options:     Options{},
	}
}
