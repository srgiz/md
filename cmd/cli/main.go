package main

import (
	"md/internal/io/cli"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := cli.Initialize()
	app.Run()
}
