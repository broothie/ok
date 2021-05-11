package task

type Arg struct {
	Parameter
	Value interface{}
}

type Args struct {
	Positional []Arg
	Keyword    map[string]Arg
}
