package cli

type Parser struct {
	version  string
	tokens   []string
	flags    []*Flag
	index    int
	taskName string
}

func NewParser(version string, tokens []string, flags ...*Flag) *Parser {
	return &Parser{
		version: version,
		tokens:  tokens,
		flags:   flags,
	}
}
