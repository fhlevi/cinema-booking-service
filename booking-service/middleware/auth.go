package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"booking-service/models"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:3001"
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		token := strings.Replace(authHeader, "Bearer ", "", 1)
		
		// Verify token with auth service
		reqBody := map[string]string{"token": token}
		jsonData, _ := json.Marshal(reqBody)
		
		resp, err := http.Post(authServiceURL+"/api/auth/verify", "application/json", bytes.NewBuffer(jsonData))
		if err != nil || resp.StatusCode != 200 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		var result struct {
			User  models.User `json:"user"`
			Valid bool        `json:"valid"`
		}
		json.NewDecoder(resp.Body).Decode(&result)

		if !result.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user", result.User)
		c.Next()
	}
}