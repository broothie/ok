package arg

type Parser struct {
	args  []string
	index int
}

func NewParser(args []string) *Parser {
	return &Parser{
		args:  args,
		index: 0,
	}
}

func (p *Parser) Advance(n int) {
	p.index += n
}

func (p *Parser) Token(index int) (Token, bool) {
	if index >= len(p.args) {
		return "", false
	}

	return Token(p.args[index]), true
}

func (p *Parser) Peek(offset int) (Token, bool) {
	return p.Token(p.index + offset)
}

func (p *Parser) Current() (Token, bool) {
	return p.Token(p.index)
}

func (p *Parser) Next() (Token, bool) {
	return p.Peek(1)
}

func (p *Parser) IsExhausted() bool {
	_, isNotExhausted := p.Current()
	return !isNotExhausted
}

func (p *Parser) HasArgsLeft() bool {
	return !p.IsExhausted()
}
