package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/iasonliu/ws/internal/handlers"
)

func routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Get("/static/", http.StripPrefix("/static/", fileServer))
	return mux
}
