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

	// Create a new Gin router
	router := routes.SetupRouter()

	// Set trusted proxies
	err = router.SetTrustedProxies([]string{"127.0.0.1"}) // Replace with your trusted proxy IPs
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Run the router
	router.Run(":8080")
}
