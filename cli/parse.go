package cli

import (
	"strings"

	"github.com/bobg/errors"
	"github.com/samber/lo"
)

func (p *Parser) TaskName() string {
	return p.taskName
}

func (p *Parser) RemainingArgs() []string {
	return p.tokens[p.index:]
}

func (p *Parser) Parse() error {
	var errs []error
	for token, tokensRemaining := p.current(); tokensRemaining; token, tokensRemaining = p.current() {
		if err := p.parseToken(token); err != nil {
			errs = append(errs, errors.Wrapf(err, "parsing token %q", token))
			p.index += 1
		}

		if p.taskName != "" {
			break
		}
	}

	return errors.Join(errs...)
}

func (p *Parser) parseToken(token string) error {
	fullFlag, rhs, _ := strings.Cut(token, "=")

	flag, err := p.parseLongFlag(fullFlag)
	if err != nil {
		return errors.Wrap(err, "parsing long flag")
	}

	if flag == nil {
		flag, err = p.parseShortFlag(fullFlag)
		if err != nil {
			return errors.Wrap(err, "parsing short flag")
		}
	}

	if flag == nil {
		p.taskName = token
		p.index += 1
		return nil
	}

	if err := p.parseFlagValue(flag, rhs); err != nil {
		return errors.Wrap(err, "parsing flag value")
	}

	return nil
}

func (p *Parser) parseLongFlag(fullFlag string) (*Flag, error) {
	name, isLongFlag := strings.CutPrefix(fullFlag, "--")
	if !isLongFlag {
		return nil, nil
	}

	flag, found := lo.Find(p.flags, func(f *Flag) bool { return f.Name == name || lo.Contains(f.Aliases, name) })
	if !found {
		return nil, errors.New("invalid flag")
	}

	return flag, nil
}

func (p *Parser) parseShortFlag(fullFlag string) (*Flag, error) {
	short, isShortFlag := strings.CutPrefix(fullFlag, "-")
	if !isShortFlag {
		return nil, nil
	}

	flag, found := lo.Find(p.flags, func(f *Flag) bool { return lo.Contains(f.Shorts, rune(short[0])) })
	if !found {
		return nil, errors.New("invalid flag")
	}

	return flag, nil
}

func (p *Parser) parseFlagValue(flag *Flag, rhs string) error {
	if flag.IsBool() {
		return p.parseFlagBoolValue(flag, rhs)
	} else {
		return p.parseFlagNonBoolValue(flag, rhs)
	}
}

func (p *Parser) parseFlagBoolValue(flag *Flag, rhs string) error {
	value := "true"
	if rhs != "" {
		value = rhs
	}

	if err := flag.Parse(value); err != nil {
		return errors.Wrapf(err, "parsing bool flag %q", flag.Name)
	}

	p.index += 1
	return nil
}

func (p *Parser) parseFlagNonBoolValue(flag *Flag, rhs string) error {
	value := rhs
	advance := 1
	if rhs == "" {
		nextToken, isPresent := p.next()
		if !isPresent {
			return errors.New("missing value for flag")
		}

		value = nextToken
		advance = 2
	}

	if err := flag.Parse(value); err != nil {
		return errors.Wrapf(err, "parsing non-bool flag %q", flag.Name)
	}

	p.index += advance
	return nil
}

func (p *Parser) next() (string, bool) {
	return p.at(p.index + 1)
}

func (p *Parser) current() (string, bool) {
	return p.at(p.index)
}

func (p *Parser) at(index int) (string, bool) {
	if !p.isInBounds(index) {
		return "", false
	}

	return p.tokens[index], true
}

func (p *Parser) isInBounds(index int) bool {
	return 0 <= index && index < len(p.tokens)
}
