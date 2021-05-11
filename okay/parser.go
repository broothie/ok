package okay

import (
	"regexp"
)

var dashPrefix = regexp.MustCompile(`^-+`)

type Parser struct {
	rawArgs    []string
	argCounter int

	options          Options
	availableOptions []Option
}

func NewParser(rawArgs []string) *Parser {
	parser := &Parser{rawArgs: rawArgs}
	parser.setupOptions()
	return parser
}

func (p *Parser) current() (string, bool) {
	return p.peek(0)
}

func (p *Parser) peek(offset int) (string, bool) {
	index := p.argCounter + offset
	if index < len(p.rawArgs) {
		return p.rawArgs[index], true
	}

	return "", false
}
