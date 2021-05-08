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
	return Parser{args: args}
}
