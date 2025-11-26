package main

import (
	"md/internal/app/server"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := server.Initialize(os.Getenv("APP_DATA"))
	app.Run()
}
