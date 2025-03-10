package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIKeyConfig holds the configuration for API key authentication
type APIKeyConfig struct {
	HeaderName  string
	ValidAPIKey string
}

// APIKeyMiddleware creates a middleware function to validate API keys
func APIKeyMiddleware(config APIKeyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the API key from the header
		apiKey := c.GetHeader(config.HeaderName)

		// Trim any whitespace
		apiKey = strings.TrimSpace(apiKey)

		// Check if the API key is empty
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing API key",
			})
			c.Abort()
			return
		}

		// Compare the API key with the valid key
		if apiKey != config.ValidAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		// If the API key is valid, continue to the next handler
		c.Next()
	}
}
