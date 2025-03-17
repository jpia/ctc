package services

import (
	"ctc/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func OverrideShortcode(c *gin.Context) {
	shortcode := c.Param("shortcode")

	urlStore := models.GetURLStore()
	url, exists := urlStore.Get(shortcode)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shortcode not found"})
		return
	}

	if url.Status == models.ReleasedStatus {
		c.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("The URL was already released via the %s method on %s.", url.ReleaseMethod, url.ReleaseTimestamp.Format(time.RFC3339))})
		return
	}

	url.Status = models.ReleasedStatus
	now := time.Now()
	url.ReleaseTimestamp = now
	url.ReleaseMethod = models.OverrideReleaseMethod
	urlStore.Set(shortcode, url, models.HighUpdatePriority)

	c.JSON(http.StatusOK, gin.H{"success": "The URL has been released early."})
}

func ListURLs(c *gin.Context) {
	urlStore := models.GetURLStore()
	urlList := urlStore.GetAll()
	c.JSON(http.StatusOK, urlList)
}

func GetStats(c *gin.Context) {
	urlStore := models.GetURLStore()
	urls := urlStore.GetAll()

	totalCount := len(urls)
	pendingCount := 0
	delayedCount := 0
	releasedCount := 0

	for _, url := range urls {
		switch url.Status {
		case models.PendingStatus:
			pendingCount++
		case models.DelayedStatus:
			delayedCount++
		case models.ReleasedStatus:
			releasedCount++
		}
	}

	stats := gin.H{
		"total_urls":     totalCount,
		"pending_count":  pendingCount,
		"delayed_count":  delayedCount,
		"released_count": releasedCount,
	}

	c.JSON(http.StatusOK, stats)
}
