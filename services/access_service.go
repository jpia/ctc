package services

import (
	"ctc/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AccessURL(c *gin.Context) {
	shortcode := c.Param("shortcode")
	url, exists := models.URLStore[shortcode]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shortcode not found"})
		return
	}

	if url.Status == models.ReleasedStatus {
		c.JSON(http.StatusOK, gin.H{"long_url": url.LongURL})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "The URL is not yet available. The release date has not passed or the weather condition does not allow access."})
	}

}
