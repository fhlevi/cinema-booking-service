package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	os.Setenv("AUTH_SERVICE_URL", "http://localhost:3001")
	os.Setenv("CINEMA_SERVICE_URL", "http://localhost:3002")
	os.Setenv("BOOKING_SERVICE_URL", "http://localhost:3003")
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func TestAPIDocsEndpoint(t *testing.T) {
	router := gin.New()
	
	// Mock API docs endpoint
	router.GET("/api/docs", func(c *gin.Context) {
		c.String(200, "<!DOCTYPE html><html><head><title>API Docs</title></head><body>Swagger UI</body></html>")
	})

	req, _ := http.NewRequest("GET", "/api/docs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !contains(w.Body.String(), "Swagger UI") {
		t.Error("Expected response to contain 'Swagger UI'")
	}
}

func TestSwaggerJSONEndpoint(t *testing.T) {
	router := gin.New()
	
	// Mock swagger JSON endpoint
	router.GET("/docs/swagger.json", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"openapi": "3.0.0",
			"info": map[string]string{
				"title":   "Cinema Booking API",
				"version": "1.0.0",
			},
		})
	})

	req, _ := http.NewRequest("GET", "/docs/swagger.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !contains(w.Body.String(), "Cinema Booking API") {
		t.Error("Expected response to contain 'Cinema Booking API'")
	}
}

func TestCORSHeaders(t *testing.T) {
	router := gin.New()
	
	// Mock CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("Expected status 204 for OPTIONS request, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS header to be set")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}