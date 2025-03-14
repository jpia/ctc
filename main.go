package main

import (
	"ctc/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := routes.SetupRouter()
	router.Run(":8080")
}
