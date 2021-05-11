package tool

import "fmt"

func Warn(toolName, format string, v ...interface{}) {
	fmt.Printf("[%s] %s\n", toolName, fmt.Sprintf(format, v...))
}
