package task

type Arg struct {
	Parameter Parameter
	Value     interface{}
}

type Args struct {
	Forwards   []string
	Positional []Arg
	Keyword    map[string]Arg
}
