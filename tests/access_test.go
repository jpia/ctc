package tests

import (
	"ctc/models"
	"ctc/routes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAccessURL(t *testing.T) {
	// Set up environment variables
	os.Setenv("USER_KEY", "user_key")
	os.Setenv("ADMIN_KEY", "admin_key")

	// Set Gin to release mode to disable debug output
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := routes.SetupRouter()

	// Add a test URL to the store
	shortcode := "test1234"
	urlStore := models.GetURLStore()
	urlStore.Reset()
	urlStore.Set(shortcode, models.URL{
		LongURL:     "https://example.com",
		ReleaseDate: time.Now().Add(24 * time.Hour),
		Shortcode:   shortcode,
		Status:      models.PendingStatus,
	}, models.HighUpdatePriority)

	// Test case: Shortcode exists but is not ready
	req, _ := http.NewRequest("GET", "/access/"+shortcode, nil)
	req.Header.Set("X-API-Key", "user_key")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "The URL is not yet available. The release date has not passed or the weather condition does not allow access.", response["error"])

	// Test case: Shortcode exists and is released
	urlStore.Set(shortcode, models.URL{
		LongURL:          "https://example.com",
		ReleaseDate:      time.Now().Add(-24 * time.Hour),
		Shortcode:        shortcode,
		Status:           models.ReleasedStatus,
		ReleaseMethod:    models.StandardReleaseMethod,
		ReleaseTimestamp: time.Now(),
	}, models.HighUpdatePriority)

	// sleep for 1 millisecond to give the URL time to be updated
	time.Sleep(1 * time.Millisecond)

	req, _ = http.NewRequest("GET", "/access/"+shortcode, nil)
	req.Header.Set("X-API-Key", "user_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	// print the response body
	fmt.Println(rr.Body.String())
	assert.Equal(t, http.StatusOK, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", response["long_url"])

	// Test case: Shortcode does not exist
	req, _ = http.NewRequest("GET", "/access/nonexistent", nil)
	req.Header.Set("X-API-Key", "user_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Shortcode not found", response["error"])
}
