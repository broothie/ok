package main

import (
	"os"

	"github.com/broothie/now/runner"
)

func main() {
	runner.Run(os.Args[1:])
}
