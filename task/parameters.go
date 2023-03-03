package task

import (
	"strings"

	"github.com/samber/lo"
)

type Parameters []Parameter

func (p Parameters) String() string {
	return strings.Join(lo.Map(p, func(param Parameter, _ int) string { return param.String() }), "  ")
}
