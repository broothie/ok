package ok

import (
	"github.com/broothie/ok/tool"
	"github.com/broothie/ok/tool/npm"
	"github.com/broothie/ok/tool/ruby"
)

func Tools() []tool.Tool {
	return []tool.Tool{
		npm.Tool{},
		ruby.Tool{},
	}
}
