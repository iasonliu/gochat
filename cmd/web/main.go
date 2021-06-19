package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/iasonliu/env"
	"github.com/iasonliu/ws/internal/handlers"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":8080", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")

func main() {
	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "GoChat",
			Level: hclog.LevelFromString(*logLevel),
		},
	)

	sm := mux.NewRouter()
	sm.HandleFunc("/ws", handlers.WsEndpoint)
	sm.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// create a logger for the server from the default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	l.Info("Starting channel listener...")
	go handlers.ListenToWsChannel()

	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Info("Starting server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
