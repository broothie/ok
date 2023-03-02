package task

import "github.com/samber/lo"

type Arguments []Argument

func (a Arguments) Required() Arguments {
	return a.Filter(func(argument Argument, _ int) bool { return argument.IsRequired() })
}

func (a Arguments) Optional() Arguments {
	return a.Filter(func(argument Argument, _ int) bool { return argument.IsOptional() })
}

func (a Arguments) Filter(predicate func(Argument, int) bool) Arguments {
	return lo.Filter(a, predicate)
}

func (a Arguments) Find(predicate func(Argument) bool) (Argument, bool) {
	return lo.Find(a, predicate)
}
