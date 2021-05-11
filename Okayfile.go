//+build okay
package main

import (
	"fmt"

	"github.com/bmatcuk/doublestar"
)

func glob(pattern string) {
	fmt.Println(doublestar.Glob(pattern))
}
