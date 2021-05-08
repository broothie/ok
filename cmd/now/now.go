package main

import (
	"os"

	"github.com/broothie/now/now"
)

func main() {
	now.Run(os.Args[1:])
}
