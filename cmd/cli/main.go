package main

import (
	_ "md/internal/infr" // init
	"md/internal/io/cli"
	"os"
)

func main() {
	if err := cli.Initialize(os.Getenv("APP_DATA")).Run(os.Args); err != nil {
		os.Exit(1)
	}
}
