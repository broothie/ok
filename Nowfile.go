package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/broothie/now/param"
)

func get(url string) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	if _, err := bufio.NewReader(res.Body).WriteTo(os.Stdout); err != nil {
		panic(err)
	}
}

func types() {
	fmt.Println(param.Untyped, param.Bool, param.Int, param.String)
}
