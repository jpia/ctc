package services

import (
	"ctc/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	shortcode := c.Param("shortcode")
	ctc, exists := models.CTCStore[shortcode]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shortcode not found"})
		return
	}

	response := gin.H{
		"status":       ctc.Status,
		"release_date": ctc.ReleaseDate,
	}

	c.JSON(http.StatusOK, response)
}
