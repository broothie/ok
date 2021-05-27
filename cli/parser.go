package cli

import (
	"regexp"

	"github.com/broothie/ok/ok"
)

var dashPrefix = regexp.MustCompile(`^-+`)

type Parser struct {
	Args       []string
	argCounter int
	options    ok.Options
}

func NewParser(args []string) (*Parser, error) {
	options, err := ok.NewOptionsFromEnvironment()
	if err != nil {
		return nil, err
	}

	return &Parser{Args: args, options: options}, nil
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
