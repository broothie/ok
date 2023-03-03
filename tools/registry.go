package tools

import (
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/tools/golang"
	"github.com/broothie/ok/tools/make"
	"github.com/broothie/ok/tools/npm"
	"github.com/broothie/ok/tools/python"
	"github.com/broothie/ok/tools/ruby"
)

func Registry() []tool.NewFunc {
	return []tool.NewFunc{
		golang.New,
		make.New,
		npm.New,
		python.New,
		ruby.New,
	}
}
