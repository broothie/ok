package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	Ok    = log.New(os.Stdout, "[ok] ", 0)
	Debug = log.New(os.Stdout, "[ok.debug] ", 0)
)

func Tool(toolName string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[ok %s]", toolName), 0)
}
