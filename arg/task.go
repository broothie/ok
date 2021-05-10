package arg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/broothie/now/param"
)

var dashStripper = regexp.MustCompile(`^-+`)

type Arg struct {
	Value interface{}
	Param param.Param
}

type Args struct {
	Positional []Arg
	Keyword    map[string]Arg
}

func newTaskArgs() Args {
	return Args{Keyword: make(map[string]Arg)}
}

func (p *Parser) ParseTaskArgs(params param.Params) error {
	p.params = params
	for p.argCounter < len(p.args) {
		if err := p.processTaskArg(); err != nil {
			return err
		}
	}

	return p.checkNumberOfArgs()
}

func (p *Parser) processTaskArg() error {
	arg, _ := p.current()
	if !strings.HasPrefix(arg, "-") {
		return p.processPositionalTaskArg()
	}

	return p.processKeywordTaskArg()
}

func (p *Parser) processPositionalTaskArg() error {
	arg, _ := p.current()
	param, ok := p.params.PositionalAt(p.positionalCount())
	if !ok {
		return fmt.Errorf("too many positional args, expected max of %d", len(p.params.PositionalRequired)+len(p.params.PositionalOptional))
	}

	castArg, err := param.Type.CastString(arg)
	if err != nil {
		return err
	}

	p.Args.Positional = append(p.Args.Positional, Arg{Value: castArg, Param: param})
	p.argCounter++
	return nil
}

func (p *Parser) processKeywordTaskArg() error {
	arg, _ := p.current()
	dashlessArg := dashStripper.ReplaceAllString(arg, "")
	parameter, found := p.findKeywordParam(dashlessArg)
	if !found && len(dashlessArg) == 1 {
		parameter, found = p.findKeywordPrefixParam(rune(dashlessArg[0]))
		if !found {
			return fmt.Errorf("no keyword param found with name '%s'", dashlessArg)
		}
	} else {
		return fmt.Errorf("no keyword param found with name '%s'", dashlessArg)
	}

	if parameter.Type == param.Bool {
		p.Args.Keyword[parameter.Name] = Arg{Value: true, Param: parameter}
		p.argCounter++
		return nil
	}

	valueArg, ok := p.peek(1)
	if !ok {
		return fmt.Errorf("no value provided for keyword arg '%s'", arg)
	}

	castValue, err := parameter.Type.CastString(valueArg)
	if err != nil {
		return err
	}

	p.Args.Keyword[parameter.Name] = Arg{Value: castValue, Param: parameter}
	p.argCounter += 2
	return nil
}

func (p *Parser) checkNumberOfArgs() error {
	if len(p.Args.Positional) < len(p.params.PositionalRequired) {
		return fmt.Errorf("missing required positional args: [%s]", p.listMissingPositionalRequiredArgs())
	}

	missingRequiredKeywords := make([]string, 0, len(p.params.KeywordRequired))
	for _, param := range p.params.KeywordRequired {
		if _, present := p.Args.Keyword[param.Name]; !present {
			missingRequiredKeywords = append(missingRequiredKeywords, param.Name)
		}
	}

	if len(missingRequiredKeywords) != 0 {
		return fmt.Errorf("missing required keyword args: [%s]", strings.Join(missingRequiredKeywords, ", "))
	}

	return nil
}

func (p *Parser) findKeywordParam(name string) (_ param.Param, found bool) {
	for _, param := range p.params.KeywordRequired {
		if param.Name == name {
			return param, true
		}
	}

	for _, param := range p.params.KeywordOptional {
		if param.Name == name {
			return param, true
		}
	}

	return param.Param{}, false
}

func (p *Parser) findKeywordPrefixParam(char rune) (_ param.Param, found bool) {
	for _, param := range p.params.KeywordRequired {
		if strings.HasPrefix(param.Name, string(char)) {
			return param, true
		}
	}

	for _, param := range p.params.KeywordOptional {
		if strings.HasPrefix(param.Name, string(char)) {
			return param, true
		}
	}

	return param.Param{}, false
}

func (p *Parser) listMissingPositionalRequiredArgs() string {
	diff := len(p.params.PositionalRequired) - len(p.Args.Positional)
	missingParams := p.params.PositionalRequired[len(p.params.PositionalRequired)-diff:]

	paramNames := make([]string, len(missingParams))
	for i, missingParam := range missingParams {
		paramNames[i] = missingParam.Name
	}

	return strings.Join(paramNames, ", ")
}
