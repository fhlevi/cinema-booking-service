package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    models.AuthRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid registration",
			requestBody: models.AuthRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing email",
			requestBody: models.AuthRequest{
				Password: "password123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Missing required fields",
		},
		{
			name: "Missing password",
			requestBody: models.AuthRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Missing required fields",
		},
		{
			name: "Missing name",
			requestBody: models.AuthRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Missing required fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/register", Register)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response["error"])
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    models.AuthRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing email",
			requestBody: models.AuthRequest{
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request",
		},
		{
			name: "Missing password",
			requestBody: models.AuthRequest{
				Email: "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/login", Login)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response["error"])
			}
		})
	}
}

func TestVerifyHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing token",
			requestBody: map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request",
		},
		{
			name: "Invalid token",
			requestBody: map[string]string{
				"token": "invalid-token",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/verify", Verify)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response["error"])
			}
		})
	}
}