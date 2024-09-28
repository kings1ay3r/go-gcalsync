package main

import (
	"gcalsync/gophers/clients/logger"
	"gcalsync/gophers/server"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create a channel to listen for interrupt signal to gracefully shut down the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go server.Serve()

	// Wait for interrupt signal
	<-stop

	logger.GetInstance().Info(nil, "Shutting down server...")

	// TODO: Implement Graceful shut down

	logger.GetInstance().Info(nil, "Server exiting")
}
