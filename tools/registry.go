package tools

import (
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/tools/bash"
	dockercompose "github.com/broothie/ok/tools/docker-compose"
	"github.com/broothie/ok/tools/golang"
	"github.com/broothie/ok/tools/make" // NOTE: Collides with `make` builtin
	"github.com/broothie/ok/tools/node"
	"github.com/broothie/ok/tools/python"
	"github.com/broothie/ok/tools/rake"
	"github.com/broothie/ok/tools/ruby"
	"github.com/broothie/ok/tools/yarn"
	"github.com/broothie/ok/tools/zsh"
)

var Registry = []tool.Tool{
	bash.Bash,
	dockercompose.Tool{},
	golang.Tool{},
	make.Make,
	node.Node,
	python.Python,
	new(rake.Tool),
	ruby.Ruby,
	yarn.Tool{},
	zsh.Zsh,
}
