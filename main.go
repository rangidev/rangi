package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rangidev/rangi/config"
	"github.com/rangidev/rangi/server"
)

func main() {
	// Config
	config := config.New()
	// Server
	server, err := server.New(config)
	if err != nil {
		config.Logger.Error("Could not create server", "error", err)
		os.Exit(1)
	}
	// Implement graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			config.Logger.Error("Could not shutdown server", "error", err)
		}
		close(idleConnsClosed)
	}()
	// Start server
	err = server.Start()
	if err != nil {
		config.Logger.Error("Could not start server", "error", err)
		os.Exit(1)
	}
	<-idleConnsClosed
}
