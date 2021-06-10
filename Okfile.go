//+build ok

package main

import "fmt"

func testSlice() {
	//args := []string{"bundle", "exec", "rails", "runner"}
	args := []string{"bundle"}
	fmt.Println("[0]", args[0])
	fmt.Println("[0:]", args[1:])
}

func testOrder() {
	fmt.Println("Okfile.go testOrder")
}
