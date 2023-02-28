package cli

import (
	"os"
	"strings"
)

type Parser struct {
	args  []string
	index int
}

func New(args []string) *Parser {
	return &Parser{
		args:  args,
		index: 0,
	}
}

func NewFromArgs() *Parser {
	return New(os.Args[1:])
}

func (p *Parser) currentIsFlag() bool {
	return p.currentIsShortFlag() || p.currentIsLongFlag()
}

func (p *Parser) currentIsShortFlag() bool {
	return strings.HasPrefix(p.current(), "-") && !p.currentIsLongFlag()
}

func (p *Parser) currentIsLongFlag() bool {
	return strings.HasPrefix(p.current(), "--")
}

func (p *Parser) currentDashless() string {
	if p.currentIsLongFlag() {
		return strings.TrimPrefix(p.current(), "--")
	} else {
		return strings.TrimPrefix(p.current(), "-")
	}
}

func (p *Parser) isDone() bool {
	return p.index >= len(p.args)
}

func (p *Parser) current() string {
	return p.args[p.index]
}

func (p *Parser) peek(n int) (string, bool) {
	index := p.index + n
	if index >= len(p.args) {
		return "", false
	}

	return p.args[index], true
}
