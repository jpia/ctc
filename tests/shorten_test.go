package tests

import (
	"bytes"
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

func TestShortenURL(t *testing.T) {
	// Set up environment variables
	os.Setenv("SHORTCODE_LENGTH", "4")
	os.Setenv("USER_KEY", "user_key")
	os.Setenv("ADMIN_KEY", "admin_key")

	// Set Gin to release mode to disable debug output
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := routes.SetupRouter()

	// Create a request body
	requestBody := map[string]interface{}{
		"long_url":     "https://example.com",
		"release_date": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
	jsonValue, _ := json.Marshal(requestBody)

	// Test with user key
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "user_key")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["shortcode"])

	fmt.Printf("User key shortcode: %s\n", response["shortcode"])

	// Test with admin key
	req, _ = http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "admin_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["shortcode"])

	fmt.Printf("Admin key shortcode: %s\n", response["shortcode"])
}

func TestShortenURLUnauthorized(t *testing.T) {
	// Set up environment variables
	os.Setenv("SHORTCODE_LENGTH", "8")
	os.Setenv("USER_KEY", "user_key")
	os.Setenv("ADMIN_KEY", "admin_key")

	// Set Gin to release mode to disable debug output
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := routes.SetupRouter()

	// Create a request body
	requestBody := map[string]interface{}{
		"long_url":     "https://example.com",
		"release_date": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
	jsonValue, _ := json.Marshal(requestBody)

	// Test without authorization header
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Test with invalid key
	req, _ = http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "invalid_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
