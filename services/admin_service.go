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

	url, exists := models.URLStore[shortcode]
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
	models.URLStore[shortcode] = url

	c.JSON(http.StatusOK, gin.H{"success": "The URL has been released early."})
}

func ListURLs(c *gin.Context) {
	var urlList []models.URL
	for _, url := range models.URLStore {
		urlList = append(urlList, url)
	}
	c.JSON(http.StatusOK, urlList)
}
