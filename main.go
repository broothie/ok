package main

import (
	_ "embed"
	"os"
	"strings"

	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/ok"
	_ "github.com/joho/godotenv/autoload"
)

//go:embed VERSION
var version string

func main() {
	if err := ok.Run(strings.TrimSpace(version), os.Args[1:]); err != nil {
		logger.Ok.Println(err)
		os.Exit(1)
	}
}
