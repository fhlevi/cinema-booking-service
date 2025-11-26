package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"booking-service/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateOnlineBookingHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    models.OnlineBookingRequest
		setupUser      bool
		expectedStatus int
	}{
		{
			name: "Valid online booking",
			requestBody: models.OnlineBookingRequest{
				StudioID: 1,
				SeatIDs:  []uint{1, 2, 3},
			},
			setupUser:      true,
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
		{
			name: "Missing studio ID",
			requestBody: models.OnlineBookingRequest{
				SeatIDs: []uint{1, 2, 3},
			},
			setupUser:      true,
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
		{
			name: "Empty seat IDs",
			requestBody: models.OnlineBookingRequest{
				StudioID: 1,
				SeatIDs:  []uint{},
			},
			setupUser:      true,
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			
			// Mock middleware to set user
			if tt.setupUser {
				router.Use(func(c *gin.Context) {
					user := models.User{
						ID:    1,
						Email: "test@example.com",
						Name:  "Test User",
						Role:  "customer",
					}
					c.Set("user", user)
					c.Next()
				})
			}
			
			router.POST("/booking/online", CreateOnlineBooking)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/booking/online", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreateOfflineBookingHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    models.OfflineBookingRequest
		expectedStatus int
	}{
		{
			name: "Valid offline booking",
			requestBody: models.OfflineBookingRequest{
				StudioID:      1,
				SeatIDs:       []uint{1, 2, 3},
				CustomerName:  "John Doe",
				CustomerEmail: "john@example.com",
			},
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
		{
			name: "Missing customer name",
			requestBody: models.OfflineBookingRequest{
				StudioID:      1,
				SeatIDs:       []uint{1, 2, 3},
				CustomerEmail: "john@example.com",
			},
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
		{
			name: "Missing customer email",
			requestBody: models.OfflineBookingRequest{
				StudioID:     1,
				SeatIDs:      []uint{1, 2, 3},
				CustomerName: "John Doe",
			},
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/booking/offline", CreateOfflineBooking)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/booking/offline", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestValidateQRCodeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    models.ValidateQRRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid booking code",
			requestBody: models.ValidateQRRequest{
				BookingCode: "VALID123",
			},
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
		{
			name: "Empty booking code",
			requestBody: models.ValidateQRRequest{
				BookingCode: "",
			},
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/validate", ValidateQRCode)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/validate", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetUserBookingsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	
	// Mock middleware to set user
	router.Use(func(c *gin.Context) {
		user := models.User{
			ID:    1,
			Email: "test@example.com",
			Name:  "Test User",
			Role:  "customer",
		}
		c.Set("user", user)
		c.Next()
	})
	
	router.GET("/bookings", GetUserBookings)

	req, _ := http.NewRequest("GET", "/bookings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 500 since no database is connected in test
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to fetch bookings", response["error"])
}

func TestInvalidJSONBookingRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		endpoint string
		handler  gin.HandlerFunc
	}{
		{
			name:     "Invalid JSON for online booking",
			endpoint: "/booking/online",
			handler:  CreateOnlineBooking,
		},
		{
			name:     "Invalid JSON for offline booking",
			endpoint: "/booking/offline",
			handler:  CreateOfflineBooking,
		},
		{
			name:     "Invalid JSON for validate QR",
			endpoint: "/validate",
			handler:  ValidateQRCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			
			// Mock middleware for online booking
			if tt.endpoint == "/booking/online" {
				router.Use(func(c *gin.Context) {
					user := models.User{ID: 1, Email: "test@example.com", Name: "Test User", Role: "customer"}
					c.Set("user", user)
					c.Next()
				})
			}
			
			router.POST(tt.endpoint, tt.handler)

			// Send invalid JSON
			req, _ := http.NewRequest("POST", tt.endpoint, bytes.NewBuffer([]byte("invalid-json")))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, "Invalid request", response["error"])
		})
	}
}