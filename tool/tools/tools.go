package tools

import (
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/tool/tools/make"
	"github.com/broothie/ok/tool/tools/npm"
)

func All() []tool.Tool {
	return []tool.Tool{
		make.NewMake(),
		npm.NewNPM(),
	}
}
