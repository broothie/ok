package main

import (
	"os"

	"github.com/broothie/now/driver"
)

func main() {
	driver.Run(os.Args[1:])
}
