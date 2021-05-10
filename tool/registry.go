package tool

import (
	"github.com/broothie/okay/tools/golang"
	"github.com/broothie/okay/tools/make" // NOTE: Collides with `make` builtin
	"github.com/broothie/okay/tools/ruby"
	"github.com/broothie/okay/tools/yarn"
)

var Registry = map[string]Tool{
	ruby.ToolName:   ruby.Ruby{},
	golang.ToolName: golang.Golang{},
	make.ToolName:   make.Make{},
	yarn.ToolName:   yarn.Yarn{},
}
