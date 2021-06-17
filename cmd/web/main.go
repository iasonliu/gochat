package main

import (
	"log"
	"net/http"
)

func main() {
	mux := routes()

	log.Println("Starting web servser on port 8080")
	_ = http.ListenAndServe(":8080", mux)
}
