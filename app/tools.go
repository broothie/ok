package app

import (
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/tool/golang"
	makepkg "github.com/broothie/ok/tool/make"
	"github.com/broothie/ok/tool/npm"
	"github.com/broothie/ok/tool/python"
	"github.com/broothie/ok/tool/ruby"
)

func Tools() []tool.Tool {
	return []tool.Tool{
		golang.Tool{},
		makepkg.Tool{},
		npm.Tool{},
		python.Tool{},
		ruby.Tool{},
	}
}
