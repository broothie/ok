package cli

func parserWithArgs(args ...string) *Parser {
	return &Parser{Args: args}
}
