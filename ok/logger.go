package ok

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "[now] ", 0)
