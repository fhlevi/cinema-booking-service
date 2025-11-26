package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	authServiceURL    string
	cinemaServiceURL  string
	bookingServiceURL string
)

func main() {
	authServiceURL = getEnv("AUTH_SERVICE_URL", "http://localhost:3001")
	cinemaServiceURL = getEnv("CINEMA_SERVICE_URL", "http://localhost:3002")
	bookingServiceURL = getEnv("BOOKING_SERVICE_URL", "http://localhost:3003")

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "api-gateway"})
	})

	// Serve Swagger UI at root docs path
	r.GET("/api/docs", func(c *gin.Context) {
		c.File("./docs/index.html")
	})
	
	// Serve swagger JSON
	r.GET("/docs/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})
	
	// Serve swagger YAML
	r.GET("/docs/swagger.yaml", func(c *gin.Context) {
		c.File("./docs/swagger.yaml")
	})

	// Proxy routes
	r.Any("/api/auth/*path", proxyHandler(authServiceURL))
	r.Any("/api/cinema/*path", proxyHandler(cinemaServiceURL))
	r.Any("/api/booking/*path", proxyHandler(bookingServiceURL))

	port := getEnv("PORT", "8080")
	log.Printf("API Gateway running on port %s", port)
	log.Printf("🚀 Swagger UI Documentation: http://localhost:%s/api/docs", port)
	log.Printf("📄 Swagger JSON: http://localhost:%s/docs/swagger.json", port)
	r.Run(":" + port)
}

func proxyHandler(targetURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build target URL
		path := c.Param("path")
		if path == "" {
			path = "/"
		}
		
		// Remove the service prefix from the path
		originalPath := c.Request.URL.Path
		if strings.HasPrefix(originalPath, "/api/auth") {
			path = strings.TrimPrefix(originalPath, "/api/auth")
		} else if strings.HasPrefix(originalPath, "/api/cinema") {
			path = strings.TrimPrefix(originalPath, "/api/cinema")
		} else if strings.HasPrefix(originalPath, "/api/booking") {
			path = strings.TrimPrefix(originalPath, "/api/booking")
		}
		
		targetPath := targetURL + "/api" + strings.TrimPrefix(originalPath, "/api")
		if c.Request.URL.RawQuery != "" {
			targetPath += "?" + c.Request.URL.RawQuery
		}

		// Create new request
		req, err := http.NewRequest(c.Request.Method, targetPath, c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create request"})
			return
		}

		// Copy headers
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Make request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": "Service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Copy response body
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}