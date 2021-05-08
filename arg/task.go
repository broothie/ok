package arg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/broothie/now/param"
)

var dashStripper = regexp.MustCompile(`^-+`)

type Task struct {
	Positional []interface{}
	Keyword    map[string]interface{}
}

func (p *Parser) ParseTaskArgs(params param.Params) error {
	p.params = params
	for p.argCounter < len(p.args) {
		if err := p.processTaskArg(); err != nil {
			return err
		}
	}

	if len(p.TaskArgs.Positional) < len(p.params.PositionalRequired) {
		return fmt.Errorf("missing required positional args [%s]", p.listMissingPositionalRequiredArgs())
	}

	if len(p.TaskArgs.Positional) > len(p.params.PositionalRequired)+len(p.params.PositionalOptional) {
		// TODO: handle too many positional args case
	}

	if len(p.TaskArgs.Keyword) < len(p.params.KeywordRequired) {
		return fmt.Errorf("missing required keyword args [%s]", p.listMissingKeywordRequiredArgs())
	}

	if len(p.TaskArgs.Keyword) > len(p.params.KeywordRequired)+len(p.params.KeywordOptional) {
		// TODO: handle too many keyword args case
	}

	return nil
}

func (p *Parser) processTaskArg() error {
	arg := p.args[p.argCounter]
	if !strings.HasPrefix(arg, "-") {
		p.TaskArgs.Positional = append(p.TaskArgs.Positional, arg)
		p.argCounter++
		return nil
	}

	dashlessArg := dashStripper.ReplaceAllString(arg, "")
	param, found := p.findKeywordParam(dashlessArg)
	if !found {
		if len(dashlessArg) == 1 {
			param, found = p.findKeywordPrefixParam(rune(dashlessArg[0]))
			if !found {
				return fmt.Errorf("no keyword param found with name '%s'", dashlessArg)
			}
		}

		return fmt.Errorf("no keyword param found with name '%s'", dashlessArg)
	}

	// TODO: check that next arg even exists
	p.TaskArgs.Keyword[param.Name] = p.args[p.argCounter+1]
	p.argCounter += 2
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
	diff := len(p.params.PositionalRequired) - len(p.TaskArgs.Positional)
	missingParams := p.params.PositionalRequired[len(p.params.PositionalRequired)-diff:]

	paramNames := make([]string, len(missingParams))
	for i, missingParam := range missingParams {
		paramNames[i] = missingParam.Name
	}

	return strings.Join(paramNames, ", ")
}

func (p *Parser) listMissingKeywordRequiredArgs() string {
	diff := len(p.params.KeywordRequired) - len(p.TaskArgs.Keyword)
	missingParams := p.params.KeywordRequired[len(p.params.KeywordRequired)-diff:]

	paramNames := make([]string, len(missingParams))
	for i, missingParam := range missingParams {
		paramNames[i] = missingParam.Name
	}

	return strings.Join(paramNames, ", ")
}
