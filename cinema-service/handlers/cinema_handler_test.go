package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cinema-service/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetStudiosHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/studios", GetStudios)

	req, _ := http.NewRequest("GET", "/studios", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 500 since no database is connected in test
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to fetch studios", response["error"])
}

func TestGetStudioSeatsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/studios/:id/seats", GetStudioSeats)

	req, _ := http.NewRequest("GET", "/studios/1/seats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 500 since no database is connected in test
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Failed to fetch seats", response["error"])
}

func TestReserveSeatsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    models.SeatReservationRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid seat reservation",
			requestBody: models.SeatReservationRequest{
				SeatIDs: []uint{1, 2, 3},
			},
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
		{
			name: "Empty seat IDs",
			requestBody: models.SeatReservationRequest{
				SeatIDs: []uint{},
			},
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/reserve", ReserveSeats)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/reserve", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestReleaseSeatsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    models.SeatReservationRequest
		expectedStatus int
	}{
		{
			name: "Valid seat release",
			requestBody: models.SeatReservationRequest{
				SeatIDs: []uint{1, 2, 3},
			},
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
		{
			name: "Empty seat IDs",
			requestBody: models.SeatReservationRequest{
				SeatIDs: []uint{},
			},
			expectedStatus: http.StatusInternalServerError, // No DB connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/release", ReleaseSeats)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/release", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestInvalidJSONRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/reserve", ReserveSeats)

	// Send invalid JSON
	req, _ := http.NewRequest("POST", "/reserve", bytes.NewBuffer([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid request", response["error"])
}