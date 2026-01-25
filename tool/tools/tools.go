package tools

import (
	"github.com/broothie/ok"
	"github.com/broothie/ok/tool/tools/just"
	"github.com/broothie/ok/tool/tools/make"
	"github.com/broothie/ok/tool/tools/npm"
	"github.com/broothie/ok/tool/tools/rake"
	"github.com/broothie/ok/tool/tools/task"
)

func All() []ok.Tool {
	return []ok.Tool{
		just.New(),
		make.New(),
		npm.New(),
		rake.New(),
		task.New(),
	}
}
