package main

import (
	"md/internal/io/cli"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := cli.Initialize(os.Getenv("APP_DATA"))
	app.Run()
}
