package task

import (
	"strings"

	"github.com/samber/lo"
)

type Parameters []Parameter

func (p Parameters) String() string {
	return strings.Join(lo.Map(p, func(param Parameter, _ int) string { return param.String() }), "  ")
}

func (p Parameters) IsSplat() bool {
	if len(p) != 1 {
		return false
	}

	return p[len(p)-1].IsSplat()
}
