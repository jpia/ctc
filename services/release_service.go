package services

import (
	"ctc/models"
	"os"
	"strconv"
	"time"
)

func ReleasePendingURLs() {
	weatherInstance := models.GetWeatherStatusInstance()
	today := time.Now().Truncate(24 * time.Hour)

	for shortcode, url := range models.URLStore {
		if (url.Status == models.PendingStatus || url.Status == models.DelayedStatus) && today.After(url.ReleaseDate.Truncate(24*time.Hour)) {
			if models.IsValidForStandardRelease(weatherInstance.Status) {
				url.Status = models.ReleasedStatus
				url.ReleaseMethod = models.StandardReleaseMethod
				url.ReleaseTimestamp = time.Now()
				models.URLStore[shortcode] = url
				DebugLog("Shortcode %s is now released due to valid weather.\n", shortcode)
			} else {
				url.Status = models.DelayedStatus
				url.ReleaseDate = today.Add(24 * time.Hour)
				url.Delays++
				models.URLStore[shortcode] = url
				DebugLog("Shortcode %s release delayed due to invalid weather.\n", shortcode)
			}
		}
	}
}

func StartReleaseService() {
	// Perform UpdateWeatherStatus instantly
	UpdateWeatherStatus()

	// Get the ticker interval from the environment variable or use the default value
	intervalStr := os.Getenv("RELEASE_TICKER_INTERVAL")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval <= 0 {
		interval = 3600 // Default to 1 hour in seconds
	}

	// Start the ticker
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		UpdateWeatherStatus()
		ReleasePendingURLs()
	}
}
