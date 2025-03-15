package main

import (
	"ctc/routes"
	"ctc/services"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize logging
	services.InitLogging()

	// Start the release service routine
	go services.StartReleaseService()

	router := routes.SetupRouter()
	router.Run(":8080")
}
