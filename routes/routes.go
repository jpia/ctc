package routes

import (
	"ctc/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func authMiddleware(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-API-Key header required"})
		c.Abort()
		return
	}

	validAPIKey := os.Getenv("USER_KEY")
	validAdminKey := os.Getenv("ADMIN_KEY")

	if apiKey != validAPIKey && apiKey != validAdminKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		c.Abort()
		return
	}

	c.Set("isAdmin", apiKey == validAdminKey)
	c.Next()
}

func adminMiddleware(c *gin.Context) {
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		c.Abort()
		return
	}
	c.Next()
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(authMiddleware)

	api := router.Group("/")
	{
		api.POST("/shorten", services.ShortenURL)
		api.GET("/status/:shortcode", services.GetStatus)
		api.GET("/access/:shortcode", services.AccessURL)
	}

	admin := router.Group("/admin")
	admin.Use(adminMiddleware)
	{
		admin.POST("/override/:shortcode", services.OverrideShortcode)
		admin.GET("/list", services.ListCTCs)
	}

	return router
}
