package services

import (
	"ctc/models"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ShortenURL(c *gin.Context) {
	var request struct {
		LongURL     string    `json:"long_url" binding:"required"`
		ReleaseDate time.Time `json:"release_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortcodeLength, err := strconv.Atoi(os.Getenv("SHORTCODE_LENGTH"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid shortcode length"})
		return
	}

	shortcode := generateShortcode(shortcodeLength)
	url := models.URL{
		LongURL:          request.LongURL,
		ReleaseDate:      request.ReleaseDate,
		ReleaseDateOrig:  request.ReleaseDate,
		Shortcode:        shortcode,
		Status:           models.PendingStatus,
		Delays:           0,
		ReleaseMethod:    "",
		ReleaseTimestamp: time.Time{},
	}

	models.URLStore[shortcode] = url

	c.JSON(http.StatusOK, gin.H{"shortcode": shortcode})
}

func generateShortcode(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	for {
		b := make([]byte, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		shortcode := string(b)
		if _, exists := models.URLStore[shortcode]; !exists {
			return shortcode
		}
	}
}
