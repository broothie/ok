package cli

import (
	"regexp"
)

var dashPrefix = regexp.MustCompile(`^-+`)

type Parser struct {
	Args       []string
	argCounter int
	options    Options
}

func NewParser(args []string) (*Parser, error) {
	return &Parser{Args: args, options: Options{Config: ReadInConfig()}}, nil
}

func (p *Parser) current() (string, bool) {
	return p.peek(0)
}

func (p *Parser) peek(offset int) (string, bool) {
	index := p.argCounter + offset
	if index < len(p.Args) {
		return p.Args[index], true
	}

	return "", false
}
