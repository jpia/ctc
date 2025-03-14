package main

import (
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

	// Create a ticker and start the UpdateStoreByWeather routine
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			<-ticker.C
			services.UpdateStoreByWeather()
		}
	}()

	router := routes.SetupRouter()
	router.Run(":8080")
}
