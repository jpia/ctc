package main

import (
	"ctc/logger"
	"ctc/routes"
	"ctc/services"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Set the timezone to match New York
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
	time.Local = loc

	// Initialize logging
	logger.InitLogging()

	// Start the release service routine
	go services.StartReleaseService()

	router := routes.SetupRouter()
	router.Run(":8080")
}
