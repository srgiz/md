package main

import (
	_ "md/internal/infr" // init
	"md/internal/io/cli"
	"os"
)

func main() {
	cli.Initialize(os.Getenv("APP_DATA")).Run()
}
