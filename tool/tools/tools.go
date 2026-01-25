package tools

import (
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/tool/tools/just"
	"github.com/broothie/ok/tool/tools/make"
	"github.com/broothie/ok/tool/tools/npm"
	"github.com/broothie/ok/tool/tools/rake"
)

func All() []tool.Tool {
	return []tool.Tool{
		just.New(),
		make.New(),
		npm.New(),
		rake.New(),
	}
}
