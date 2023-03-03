package tools

import (
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/tools/bash"
	"github.com/broothie/ok/tools/dockercompose"
	"github.com/broothie/ok/tools/golang"
	"github.com/broothie/ok/tools/make"
	"github.com/broothie/ok/tools/node"
	"github.com/broothie/ok/tools/npm"
	"github.com/broothie/ok/tools/python"
	"github.com/broothie/ok/tools/ruby"
)

func Registry() []tool.NewFunc {
	return []tool.NewFunc{
		bash.New,
		dockercompose.New,
		golang.New,
		make.New,
		node.New,
		npm.New,
		python.New,
		ruby.New,
	}
}
