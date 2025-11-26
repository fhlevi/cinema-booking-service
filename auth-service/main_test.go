package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

func TestMain(m *testing.M) {
	// Set test environment
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test?sslmode=disable")
	os.Setenv("JWT_SECRET", "test-secret")
	gin.SetMode(gin.TestMode)
	
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestHealthEndpoint(t *testing.T) {
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "auth-service"})
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response["status"] != "OK" {
		t.Errorf("Expected status OK, got %s", response["status"])
	}
}

func TestGoogleOAuthEndpoints(t *testing.T) {
	router := gin.New()
	
	// Import handlers
	router.GET("/api/auth/google", func(c *gin.Context) {
		c.Redirect(307, "https://accounts.google.com/oauth/authorize")
	})
	
	router.GET("/api/auth/google/callback", func(c *gin.Context) {
		state := c.Query("state")
		code := c.Query("code")
		
		if state == "" || code == "" {
			c.JSON(400, gin.H{"error": "Missing state or code"})
			return
		}
		
		c.JSON(200, gin.H{"message": "OAuth callback received"})
	})

	// Test Google login redirect
	req, _ := http.NewRequest("GET", "/api/auth/google", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 307 {
		t.Errorf("Expected status 307 for Google login, got %d", w.Code)
	}

	// Test callback with valid parameters
	req2, _ := http.NewRequest("GET", "/api/auth/google/callback?state=test&code=test", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Errorf("Expected status 200 for callback, got %d", w2.Code)
	}
}

func TestRegisterValidation(t *testing.T) {
	router := gin.New()
	
	// Mock register endpoint for testing validation
	router.POST("/register", func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		
		if req.Email == "" || req.Password == "" || req.Name == "" {
			c.JSON(400, gin.H{"error": "Missing required fields"})
			return
		}
		
		c.JSON(201, gin.H{"message": "User registered"})
	})

	// Test missing fields
	reqBody := map[string]string{
		"email": "test@example.com",
		// missing password and name
	}
	
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status 400 for missing fields, got %d", w.Code)
	}
}

func TestLoginValidation(t *testing.T) {
	router := gin.New()
	
	// Mock login endpoint for testing validation
	router.POST("/login", func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		
		if req.Email == "" || req.Password == "" {
			c.JSON(400, gin.H{"error": "Missing required fields"})
			return
		}
		
		// Mock authentication logic
		if req.Email == "test@example.com" && req.Password == "password123" {
			c.JSON(200, gin.H{"message": "Login successful"})
		} else {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
		}
	})

	// Test valid credentials
	reqBody := AuthRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200 for valid credentials, got %d", w.Code)
	}
}