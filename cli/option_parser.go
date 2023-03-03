package cli

type optionParser struct {
	args  []string
	index int
}

func newOptionParser(args []string) *optionParser {
	return &optionParser{
		args:  args,
		index: 0,
	}
}

func (p *optionParser) advance(n int) {
	p.index += n
}

func (p *optionParser) token(index int) (token, bool) {
	if index >= len(p.args) {
		return "", false
	}

	return token(p.args[index]), true
}

func (p *optionParser) peek(offset int) (token, bool) {
	return p.token(p.index + offset)
}

func (p *optionParser) current() (token, bool) {
	return p.token(p.index)
}

func (p *optionParser) next() (token, bool) {
	return p.peek(1)
}

func (p *optionParser) isExhausted() bool {
	_, isNotExhausted := p.current()
	return !isNotExhausted
}

func (p *optionParser) hasArgsLeft() bool {
	return !p.isExhausted()
}
