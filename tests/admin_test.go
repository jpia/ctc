package tests

import (
	"ctc/models"
	"ctc/routes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOverrideShortcode(t *testing.T) {
	// Set up environment variables
	os.Setenv("USER_KEY", "user_key")
	os.Setenv("ADMIN_KEY", "admin_key")

	// Set Gin to release mode to disable debug output
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := routes.SetupRouter()

	// Add a test URL to the store
	// generate a random 6 digit shortcode
	shortcode := "test1234"
	urlStore := models.GetURLStore()
	urlStore.Reset()
	urlStore.Set(shortcode, models.URL{
		LongURL:     "https://example.com",
		ReleaseDate: time.Now().Add(24 * time.Hour),
		Shortcode:   shortcode,
		Status:      models.PendingStatus,
	}, models.HighUpdatePriority)

	time.Sleep(1 * time.Millisecond)

	// Test case: Shortcode exists and overriden
	req, _ := http.NewRequest("POST", "/admin/override/"+shortcode, nil)
	req.Header.Set("X-API-Key", "admin_key")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["success"], "The URL has been released early.")

	// Test case: Shortcode exists and is already released
	shortcode = "test5678"
	urlStore.Set(shortcode, models.URL{
		LongURL:          "https://example.com",
		ReleaseDate:      time.Now().Add(-24 * time.Hour), // Release date in the past
		Shortcode:        shortcode,
		Status:           models.ReleasedStatus,
		ReleaseMethod:    models.StandardReleaseMethod,
		ReleaseTimestamp: time.Now().Add(-24 * time.Hour),
	}, models.HighUpdatePriority)
	time.Sleep(1 * time.Millisecond)

	req, _ = http.NewRequest("POST", "/admin/override/"+shortcode, nil)
	req.Header.Set("X-API-Key", "admin_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["success"], models.StandardReleaseMethod)

	// Test case: Shortcode does not exist
	req, _ = http.NewRequest("POST", "/admin/override/nonexistent", nil)
	req.Header.Set("X-API-Key", "admin_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Shortcode not found", response["error"])
}

func TestListURLs(t *testing.T) {
	// Set up environment variables
	os.Setenv("USER_KEY", "user_key")
	os.Setenv("ADMIN_KEY", "admin_key")

	// Set Gin to release mode to disable debug output
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := routes.SetupRouter()

	// Add test URLs to the store
	urlStore := models.GetURLStore()
	urlStore.Reset()
	urlStore.Set("test1234", models.URL{
		LongURL:     "https://example.com",
		ReleaseDate: time.Now().Add(24 * time.Hour),
		Shortcode:   "test1234",
		Status:      models.PendingStatus,
	}, models.HighUpdatePriority)
	urlStore.Set("test5678", models.URL{
		LongURL:     "https://example2.com",
		ReleaseDate: time.Now().Add(48 * time.Hour),
		Shortcode:   "test5678",
		Status:      models.PendingStatus,
	}, models.HighUpdatePriority)

	// Test case: List all URLs
	req, _ := http.NewRequest("GET", "/admin/list", nil)
	req.Header.Set("X-API-Key", "admin_key")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.URL
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}
