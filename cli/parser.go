package cli

type parser struct {
	args  []string
	index int
}

func newParser(args []string) *parser {
	return &parser{
		args:  args,
		index: 0,
	}
}

func (p *parser) advance(n int) {
	p.index += n
}

func (p *parser) token(index int) (token, bool) {
	if index >= len(p.args) {
		return "", false
	}

	return token(p.args[index]), true
}

func (p *parser) peek(offset int) (token, bool) {
	return p.token(p.index + offset)
}

func (p *parser) current() (token, bool) {
	return p.token(p.index)
}

func (p *parser) next() (token, bool) {
	return p.peek(1)
}

func (p *parser) isExhausted() bool {
	_, isNotExhausted := p.current()
	return !isNotExhausted
}

func (p *parser) hasArgsLeft() bool {
	return !p.isExhausted()
}
