package ok

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "[ok] ", 0)
