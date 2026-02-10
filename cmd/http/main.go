package main

import (
	"log"
	"net/http"
)

func main() {
	// todo.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
