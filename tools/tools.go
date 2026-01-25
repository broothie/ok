package tools

import (
	"github.com/broothie/ok"
	"github.com/broothie/ok/tools/just"
	"github.com/broothie/ok/tools/make"
	"github.com/broothie/ok/tools/packagejson"
	"github.com/broothie/ok/tools/rake"
	"github.com/broothie/ok/tools/task"
)

func All() []ok.Tool {
	return []ok.Tool{
		just.New(),
		make.New(),
		packagejson.NewNPM(),
		packagejson.NewYarn(),
		rake.New(),
		task.New(),
	}
}
