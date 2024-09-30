package main

import (
	"context"
	"fmt"
	"gcalsync/gophers/core"
	"gcalsync/gophers/dao"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Running scheduled renewal of watches...")
	err := dao.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	coreInstance, err := core.New()
	if err != nil {
		log.Fatalf("Error creating core instance: %v", err)
	}
	coreInstance.RenewExpiringWatches(context.Background())
	fmt.Println("Exiting ...")

}
