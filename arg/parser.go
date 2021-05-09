package arg

import "github.com/broothie/now/param"

type Parser struct {
	NowArgs  Now
	TaskArgs Task

	args       []string
	argCounter int
	params     param.Params
}

func NewParser(args []string) Parser {
	return Parser{
		args: args,
		// Zero value works for everything else
	}
}

func (p *Parser) current() (string, bool) {
	return p.peek(0)
}

func (p *Parser) peek(offset int) (string, bool) {
	index := p.argCounter + offset
	if index < len(p.args) {
		return p.args[index], true
	}

	return "", false
}

func (p *Parser) positionalCount() int {
	return len(p.TaskArgs.Positional)
}
