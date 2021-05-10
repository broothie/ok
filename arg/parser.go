package arg

import "github.com/broothie/okay/param"

type Parser struct {
	Options Options
	Args    Args

	args       []string
	argCounter int
	params     param.Params
}

func NewParser(args []string) Parser {
	return Parser{
		Args: newTaskArgs(),
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
	return len(p.Args.Positional)
}
