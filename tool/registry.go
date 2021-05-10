package tool

import (
	"github.com/broothie/now/tools/golang"
	"github.com/broothie/now/tools/make" // NOTE: Collides with `make` builtin
	"github.com/broothie/now/tools/ruby"
	"github.com/broothie/now/tools/yarn"
)

var Registry = map[string]Tool{
	ruby.ToolName:   ruby.Ruby{},
	golang.ToolName: golang.Golang{},
	make.ToolName:   make.Make{},
	yarn.ToolName:   yarn.Yarn{},
}
