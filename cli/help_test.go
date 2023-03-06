package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/broothie/ok"
	"github.com/stretchr/testify/assert"
)

func TestCLI_PrintHelp(t *testing.T) {
	expected := fmt.Sprintf(`ok %s

Usage:
  ok [OPTIONS] <TASK> [TASK ARGS]

Options:
  -h  --help          Show help.
  -V  --version       Show version.
      --tools         List available tools.
      --tool [TOOL]   Configure a tool. Can be used multiple times.
      --init <TOOL>   Initialize a tool.
  -w  --watch <GLOB>  Glob pattern of files to watch. Can be used multiple times.
`,
		ok.Version(),
	)

	var buf bytes.Buffer
	assert.NoError(t, (&CLI{flags: flags()}).PrintHelp(&buf))
	assert.Equal(t, expected, buf.String())
}
