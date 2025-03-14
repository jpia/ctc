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

func TestAccessURL(t *testing.T) {
	// Set up environment variables
	os.Setenv("USER_KEY", "user_key")
	os.Setenv("ADMIN_KEY", "admin_key")

	// Set Gin to release mode to disable debug output
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := routes.SetupRouter()

	// Add a test CTC to the store
	shortcode := "test1234"
	models.CTCStore[shortcode] = models.CTC{
		LongURL:     "https://example.com",
		ReleaseDate: time.Now().Add(24 * time.Hour),
		Shortcode:   shortcode,
		Status:      models.Pending,
	}

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

	// Test case: Shortcode exists and is ready but release date is in the future
	models.CTCStore[shortcode] = models.CTC{
		LongURL:     "https://example.com",
		ReleaseDate: time.Now().Add(24 * time.Hour),
		Shortcode:   shortcode,
		Status:      models.Ready,
	}

	req, _ = http.NewRequest("GET", "/access/"+shortcode, nil)
	req.Header.Set("X-API-Key", "user_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "The URL is not yet available. The release date has not passed or the weather condition does not allow access.", response["error"])

	// Test case: Shortcode exists and is ready and release date is in the past
	models.CTCStore[shortcode] = models.CTC{
		LongURL:     "https://example.com",
		ReleaseDate: time.Now().Add(-24 * time.Hour),
		Shortcode:   shortcode,
		Status:      models.Ready,
	}

	req, _ = http.NewRequest("GET", "/access/"+shortcode, nil)
	req.Header.Set("X-API-Key", "user_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

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
