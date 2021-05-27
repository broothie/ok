package cli

import (
	"regexp"
)

var dashPrefix = regexp.MustCompile(`^-+`)

type Parser struct {
	Args       []string
	argCounter int
	config     Options
}

func NewParser(args []string) (*Parser, error) {
	return &Parser{Args: args, config: ReadInNewConfig()}, nil
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
