package main

import (
	"os"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/ok"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if err := ok.Run(os.Args[1:]); err != nil {
		logger.Ok.Println(err)
		os.Exit(1)
	}
}
