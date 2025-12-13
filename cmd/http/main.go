package main

import (
	"md/internal/io/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := http.Initialize(os.Getenv("APP_DATA"))
	app.Run()
}
