package okay

func parserWithArgs(args ...string) *Parser {
	return NewParser(args)
}
