//+build ok

package main

import (
	"fmt"
	"strconv"

	"github.com/bmatcuk/doublestar"
)

// Try out doublestar
func glob(pattern string) {
	fmt.Println(doublestar.Glob(pattern))
}

func float(s string, f float64) {
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println(res == f)
}
