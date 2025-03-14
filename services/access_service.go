package services

import (
	"ctc/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AccessURL(c *gin.Context) {
	shortcode := c.Param("shortcode")
	ctc, exists := models.CTCStore[shortcode]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shortcode not found"})
		return
	}

	if ctc.Status == models.Ready && time.Now().After(ctc.ReleaseDate) {
		c.JSON(http.StatusOK, gin.H{"long_url": ctc.LongURL})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "The URL is not yet available. The release date has not passed or the weather condition does not allow access."})
	}
}
