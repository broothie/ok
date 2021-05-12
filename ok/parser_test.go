package ok

func parserWithArgs(args ...string) *Parser {
	return &Parser{Args: args}
}
