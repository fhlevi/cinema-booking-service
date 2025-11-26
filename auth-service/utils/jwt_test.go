package utils

import (
	"os"
	"testing"
	"time"

	"auth-service/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret")
	jwtSecret = []byte("test-secret")

	user := models.User{
		ID:    1,
		Email: "test@example.com",
		Role:  "customer",
	}

	token, err := GenerateToken(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token structure
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims := parsedToken.Claims.(jwt.MapClaims)
	assert.Equal(t, float64(1), claims["userId"])
	assert.Equal(t, "test@example.com", claims["email"])
	assert.Equal(t, "customer", claims["role"])
}

func TestValidateToken(t *testing.T) {
	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret")
	jwtSecret = []byte("test-secret")

	tests := []struct {
		name        string
		token       string
		expectValid bool
	}{
		{
			name:        "Invalid token format",
			token:       "invalid-token",
			expectValid: false,
		},
		{
			name:        "Empty token",
			token:       "",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ValidateToken(tt.token)

			if tt.expectValid {
				assert.NoError(t, err)
				assert.True(t, token.Valid)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidTokenGeneration(t *testing.T) {
	// Set test JWT secret
	os.Setenv("JWT_SECRET", "test-secret")
	jwtSecret = []byte("test-secret")

	user := models.User{
		ID:    1,
		Email: "test@example.com",
		Role:  "customer",
	}

	// Generate token
	tokenString, err := GenerateToken(user)
	assert.NoError(t, err)

	// Validate the generated token
	token, err := ValidateToken(tokenString)
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims := token.Claims.(jwt.MapClaims)
	assert.Equal(t, float64(1), claims["userId"])
	assert.Equal(t, "test@example.com", claims["email"])
	assert.Equal(t, "customer", claims["role"])

	// Check expiration
	exp := claims["exp"].(float64)
	assert.True(t, exp > float64(time.Now().Unix()))
}