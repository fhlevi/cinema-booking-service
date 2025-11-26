package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGoogleLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/auth/google", GoogleLogin)

	req, _ := http.NewRequest("GET", "/auth/google", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 307 {
		t.Errorf("Expected status 307 (redirect), got %d", w.Code)
	}

	location := w.Header().Get("Location")
	if !strings.Contains(location, "accounts.google.com") {
		t.Error("Expected redirect to Google OAuth")
	}
}

func TestGoogleCallbackInvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/auth/google/callback", GoogleCallback)

	// Test with invalid state
	req, _ := http.NewRequest("GET", "/auth/google/callback?state=invalid&code=test", nil)
	req.AddCookie(&http.Cookie{Name: "oauthstate", Value: "valid"})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status 400 for invalid state, got %d", w.Code)
	}
}

func TestGoogleCallbackMissingCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/auth/google/callback", GoogleCallback)

	// Test with valid state but missing code
	req, _ := http.NewRequest("GET", "/auth/google/callback?state=valid", nil)
	req.AddCookie(&http.Cookie{Name: "oauthstate", Value: "valid"})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500 for missing code, got %d", w.Code)
	}
}

func TestGenerateStateOauthCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		state := generateStateOauthCookie(c)
		c.JSON(200, gin.H{"state": state})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check if cookie is set
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "oauthstate" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected oauthstate cookie to be set")
	}
}

func TestFindOrCreateGoogleUser(t *testing.T) {
	googleUser := GoogleUser{
		ID:    "123",
		Email: "test@gmail.com",
		Name:  "Test User",
	}

	user, err := findOrCreateGoogleUser(googleUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.Email != googleUser.Email {
		t.Errorf("Expected email %s, got %s", googleUser.Email, user.Email)
	}

	if user.Name != googleUser.Name {
		t.Errorf("Expected name %s, got %s", googleUser.Name, user.Name)
	}

	if user.Role != "customer" {
		t.Errorf("Expected role 'customer', got %s", user.Role)
	}
}
