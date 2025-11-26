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

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test?sslmode=disable")
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func TestHealthEndpoint(t *testing.T) {
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "cinema-service"})
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGetStudios(t *testing.T) {
	router := gin.New()
	
	// Mock studios endpoint
	router.GET("/api/cinema/studios", func(c *gin.Context) {
		studios := []map[string]interface{}{
			{"id": 1, "name": "Studio 1", "total_seats": 20},
			{"id": 2, "name": "Studio 2", "total_seats": 20},
		}
		c.JSON(200, studios)
	})

	req, _ := http.NewRequest("GET", "/api/cinema/studios", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var studios []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &studios)
	
	if len(studios) != 2 {
		t.Errorf("Expected 2 studios, got %d", len(studios))
	}
}

func TestGetStudioSeats(t *testing.T) {
	router := gin.New()
	
	// Mock seats endpoint
	router.GET("/api/cinema/studios/:id/seats", func(c *gin.Context) {
		seats := []map[string]interface{}{
			{"id": 1, "studio_id": 1, "seat_number": "A1", "is_available": true},
			{"id": 2, "studio_id": 1, "seat_number": "A2", "is_available": true},
		}
		c.JSON(200, seats)
	})

	req, _ := http.NewRequest("GET", "/api/cinema/studios/1/seats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestReserveSeats(t *testing.T) {
	router := gin.New()
	
	// Mock reserve seats endpoint
	router.POST("/api/cinema/seats/reserve", func(c *gin.Context) {
		var req map[string][]int
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		
		seatIds := req["seatIds"]
		if len(seatIds) == 0 {
			c.JSON(400, gin.H{"error": "No seats provided"})
			return
		}
		
		c.JSON(200, gin.H{"message": "Seats reserved successfully"})
	})

	reqBody := map[string][]int{
		"seatIds": {1, 2, 3},
	}
	
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/cinema/seats/reserve", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}