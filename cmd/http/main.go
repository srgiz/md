package main

import (
	_ "md/internal/infr" // init
	"md/internal/io/http"
	"os"
)

func main() {
	http.Initialize(os.Getenv("APP_DATA")).Run()
}
