package services

import (
	"ctc/logger"
	"ctc/models"
	"fmt"
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

	shortcode, err := generateShortcode(shortcodeLength)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate unique shortcode"})
		return
	}

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

	urlStore := models.GetURLStore()
	urlStore.Set(shortcode, url, models.HighUpdatePriority)

	c.JSON(http.StatusOK, gin.H{"shortcode": shortcode})
}

func generateShortcode(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	urlStore := models.GetURLStore()

	for attempts := 0; attempts < 10; attempts++ {
		b := make([]byte, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		shortcode := string(b)
		if _, exists := urlStore.Get(shortcode); !exists {
			return shortcode, nil
		}
		logger.DebugLog("Attempt %d: Shortcode %s already exists", attempts+1, shortcode)
	}

	return "", fmt.Errorf("failed to generate unique shortcode after 10 attempts")
}
