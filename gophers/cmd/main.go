package main

import (
	"gcalsync/gophers/server"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Start the server
	server.Serve()
}
