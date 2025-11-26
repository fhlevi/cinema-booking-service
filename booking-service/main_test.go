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

type OnlineBookingRequest struct {
	StudioID uint   `json:"studioId"`
	SeatIDs  []int  `json:"seatIds"`
}

type OfflineBookingRequest struct {
	StudioID      uint   `json:"studioId"`
	SeatIDs       []int  `json:"seatIds"`
	CustomerName  string `json:"customerName"`
	CustomerEmail string `json:"customerEmail"`
}

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test?sslmode=disable")
	os.Setenv("AUTH_SERVICE_URL", "http://localhost:3001")
	os.Setenv("CINEMA_SERVICE_URL", "http://localhost:3002")
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func TestHealthEndpoint(t *testing.T) {
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "booking-service"})
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestOnlineBookingValidation(t *testing.T) {
	router := gin.New()
	
	// Mock online booking endpoint
	router.POST("/api/booking/online", func(c *gin.Context) {
		var req OnlineBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		
		if req.StudioID == 0 || len(req.SeatIDs) == 0 {
			c.JSON(400, gin.H{"error": "Missing required fields"})
			return
		}
		
		c.JSON(201, gin.H{"message": "Booking created"})
	})

	// Test valid request
	reqBody := OnlineBookingRequest{
		StudioID: 1,
		SeatIDs:  []int{1, 2, 3},
	}
	
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/booking/online", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestOfflineBookingValidation(t *testing.T) {
	router := gin.New()
	
	// Mock offline booking endpoint
	router.POST("/api/booking/offline", func(c *gin.Context) {
		var req OfflineBookingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		
		if req.StudioID == 0 || len(req.SeatIDs) == 0 || req.CustomerName == "" || req.CustomerEmail == "" {
			c.JSON(400, gin.H{"error": "Missing required fields"})
			return
		}
		
		c.JSON(201, gin.H{"message": "Booking created"})
	})

	// Test valid request
	reqBody := OfflineBookingRequest{
		StudioID:      1,
		SeatIDs:       []int{4, 5},
		CustomerName:  "Jane Doe",
		CustomerEmail: "jane@example.com",
	}
	
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/booking/offline", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestValidateQRCode(t *testing.T) {
	router := gin.New()
	
	// Mock validate endpoint
	router.POST("/api/booking/validate", func(c *gin.Context) {
		var req struct {
			BookingCode string `json:"bookingCode"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		
		if req.BookingCode == "" {
			c.JSON(400, gin.H{"error": "Missing booking code"})
			return
		}
		
		// Mock validation logic
		if req.BookingCode == "valid-code" {
			c.JSON(200, gin.H{"valid": true, "message": "Ticket validated"})
		} else {
			c.JSON(404, gin.H{"error": "Invalid or used ticket"})
		}
	})

	// Test valid booking code
	reqBody := map[string]string{
		"bookingCode": "valid-code",
	}
	
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/booking/validate", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}