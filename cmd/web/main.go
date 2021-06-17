package main

import (
	"log"
	"net/http"

	"github.com/iasonliu/ws/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Starting web servser on port 8080")
	_ = http.ListenAndServe(":8080", mux)
}
