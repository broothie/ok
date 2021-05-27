package ok

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "[ok] ", 0)
var DebugLogger = log.New(os.Stdout, "[ok.debug] ", 0)
