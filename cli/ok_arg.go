package cli

import "fmt"

type applyFunc func(*Parser, *OkArgs) error

type okArg struct {
	Name        string
	Short       rune
	Description string
	Apply       applyFunc
}

func (a okArg) HasShort() bool {
	return a.Short != 0
}

func (a okArg) Match(argString string) bool {
	return argString == fmt.Sprintf("--%s", a.Name) || argString == fmt.Sprintf("-%c", a.Short)
}
