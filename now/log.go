package now

import "fmt"

func log(format string, v ...interface{}) {
	fmt.Printf("[now] %s\n", fmt.Sprintf(format, v...))
}
