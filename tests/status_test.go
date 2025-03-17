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

func TestGetStatus(t *testing.T) {
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

	time.Sleep(1 * time.Millisecond)

	// Test case: Shortcode exists
	req, _ := http.NewRequest("GET", "/status/"+shortcode, nil)
	req.Header.Set("X-API-Key", "user_key")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, string(models.PendingStatus), response["status"])
	assert.NotEmpty(t, response["release_date"])

	// Test case: Shortcode does not exist
	req, _ = http.NewRequest("GET", "/status/nonexistent", nil)
	req.Header.Set("X-API-Key", "user_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Shortcode not found", response["error"])
}
