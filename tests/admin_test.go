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

	// Add a test CTC to the store
	shortcode := "test1234"
	models.CTCStore[shortcode] = models.CTC{
		LongURL:     "https://example.com",
		ReleaseDate: time.Now().Add(24 * time.Hour),
		Shortcode:   shortcode,
		Status:      models.Pending,
	}

	// Test case: Shortcode exists and is not released
	req, _ := http.NewRequest("POST", "/admin/override/"+shortcode, nil)
	req.Header.Set("X-API-Key", "admin_key")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "The URL has been released early.", response["success"])

	// Test case: Shortcode exists and is already released
	models.CTCStore[shortcode] = models.CTC{
		LongURL:     "https://example.com",
		ReleaseDate: time.Now().Add(-24 * time.Hour), // Release date in the past
		Shortcode:   shortcode,
		Status:      models.Ready,
	}

	req, _ = http.NewRequest("POST", "/admin/override/"+shortcode, nil)
	req.Header.Set("X-API-Key", "admin_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "The URL has already been released.", response["success"])

	// Test case: Shortcode exists and is not released but release date is in the past
	models.CTCStore[shortcode] = models.CTC{
		LongURL:     "https://example.com",
		ReleaseDate: time.Now().Add(-24 * time.Hour), // Release date in the past
		Shortcode:   shortcode,
		Status:      models.Pending,
	}

	req, _ = http.NewRequest("POST", "/admin/override/"+shortcode, nil)
	req.Header.Set("X-API-Key", "admin_key")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "The URL has been released early.", response["success"])

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

func TestListCTCs(t *testing.T) {
	// Set up environment variables
	os.Setenv("USER_KEY", "user_key")
	os.Setenv("ADMIN_KEY", "admin_key")

	// Set Gin to release mode to disable debug output
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := routes.SetupRouter()

	// Add test CTCs to the store
	models.CTCStore["test1234"] = models.CTC{
		LongURL:     "https://example.com",
		ReleaseDate: time.Now().Add(24 * time.Hour),
		Shortcode:   "test1234",
		Status:      models.Pending,
	}
	models.CTCStore["test5678"] = models.CTC{
		LongURL:     "https://example2.com",
		ReleaseDate: time.Now().Add(48 * time.Hour),
		Shortcode:   "test5678",
		Status:      models.Pending,
	}

	// Test case: List all CTCs
	req, _ := http.NewRequest("GET", "/admin/list", nil)
	req.Header.Set("X-API-Key", "admin_key")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []models.CTC
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}
