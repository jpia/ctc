package services

import (
	"fmt"
	"log"
	"os"
	"time"
)

func InitLogging() {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatalf("Error creating log directory: %v", err)
	}

	// Set up logging to a file with a timestamp suffix
	logFileName := fmt.Sprintf("logs/app_%s.log", time.Now().Format("20060102_150405"))
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	log.SetOutput(logFile)

}

func DebugLog(format string, v ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func ErrorLog(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}
