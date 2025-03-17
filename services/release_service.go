package services

import (
	"ctc/logger"
	"ctc/models"
	"os"
	"strconv"
	"time"
)

func ReleasePendingURLs() {
	startTime := time.Now() // Start tracking time

	weatherInstance := models.GetWeatherStatusInstance()
	today := time.Now()

	const layout = "2006-01-02"
	if weatherInstance.DateChecked.Format(layout) != today.Format(layout) {
		logger.DebugLog("Current date does not match checked date, will use checked date for safety \n")
		today = weatherInstance.DateChecked
	}

	urlStore := models.GetURLStore()
	urls := urlStore.GetAll()

	eligibleCount := 0
	totalCount := len(urls)

	for _, url := range urls {
		if (url.Status == models.PendingStatus || url.Status == models.DelayedStatus) && today.After(url.ReleaseDate) {
			eligibleCount++
			if models.IsValidForStandardRelease(weatherInstance.Status) {
				url.Status = models.ReleasedStatus
				url.ReleaseMethod = models.StandardReleaseMethod
				url.ReleaseTimestamp = time.Now()
				urlStore.Set(url.Shortcode, url, models.LowUpdatePriority)
				logger.DebugLog("Shortcode %s is now released due to valid weather.\n", url.Shortcode)
			} else if models.IsValidForApiSickDayRelease(weatherInstance.Status) {
				url.Status = models.ReleasedStatus
				url.ReleaseMethod = models.ApiSickDayReleaseMethod
				url.ReleaseTimestamp = time.Now()
				urlStore.Set(url.Shortcode, url, models.LowUpdatePriority)
				logger.DebugLog("Shortcode %s is now released due to API Sick Day.\n", url.Shortcode)
			} else {
				url.Status = models.DelayedStatus
				url.ReleaseDate = today.Add(24 * time.Hour)
				url.Delays++
				urlStore.Set(url.Shortcode, url, models.LowUpdatePriority)
				logger.DebugLog("Shortcode %s release delayed due to invalid weather.\n", url.Shortcode)
			}
		}
	}

	duration := time.Since(startTime) // Calculate the duration
	logger.InfoLog("ReleasePendingURLs took %s and checked %d eligible URLs out of %d total URLs", duration, eligibleCount, totalCount)
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
