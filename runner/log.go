package runner

import "fmt"

func log(format string, v ...interface{}) {
	fmt.Printf("[okay] %s\n", fmt.Sprintf(format, v...))
}
