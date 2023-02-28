package logger

import (
	"log"
	"os"
)

var Log = log.New(os.Stdout, "[ok] ", 0)
