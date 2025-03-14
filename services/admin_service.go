package services

import (
	"ctc/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func OverrideShortcode(c *gin.Context) {
	shortcode := c.Param("shortcode")

	ctc, exists := models.CTCStore[shortcode]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shortcode not found"})
		return
	}

	if ctc.Status == models.Ready && ctc.ReleaseDate.Before(time.Now()) {
		c.JSON(http.StatusOK, gin.H{"success": "The URL has already been released."})
		return
	}

	ctc.Status = models.Ready
	if ctc.ReleaseDate.After(time.Now()) {
		ctc.ReleaseDate = time.Now()
	}
	models.CTCStore[shortcode] = ctc

	c.JSON(http.StatusOK, gin.H{"success": "The URL has been released early."})
}

func ListCTCs(c *gin.Context) {
	var ctcList []models.CTC
	for _, ctc := range models.CTCStore {
		ctcList = append(ctcList, ctc)
	}
	c.JSON(http.StatusOK, ctcList)
}
