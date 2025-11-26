package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"auth-service/models"
	"auth-service/services"
	"auth-service/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:3001/api/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func GoogleLogin(c *gin.Context) {
	oauthState := generateStateOauthCookie(c)
	url := googleOauthConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	oauthState, _ := c.Cookie("oauthstate")
	
	if c.Request.FormValue("state") != oauthState {
		c.JSON(400, gin.H{"error": "Invalid oauth state"})
		return
	}

	data, err := getUserDataFromGoogle(c.Request.FormValue("code"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user data from Google"})
		return
	}

	var googleUser GoogleUser
	if err := json.Unmarshal(data, &googleUser); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse user data"})
		return
	}

	// Check if user exists or create new user
	user, err := findOrCreateGoogleUser(googleUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process user"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(*user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"user":  user,
		"token": token,
	})
}

func generateStateOauthCookie(c *gin.Context) string {
	state := "random-state-string" // In production, use crypto/rand
	c.SetCookie("oauthstate", state, 3600, "/", "localhost", false, true)
	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents := make([]byte, 1024)
	n, err := response.Body.Read(contents)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	return contents[:n], nil
}

func findOrCreateGoogleUser(googleUser GoogleUser) (*models.User, error) {
	// Try to find existing user by email
	user, err := services.GetUserByEmail(googleUser.Email)
	if err == nil {
		return user, nil
	}

	// Create new user
	req := models.AuthRequest{
		Email:    googleUser.Email,
		Name:     googleUser.Name,
		Password: "google-oauth", // Placeholder password for OAuth users
	}

	return services.RegisterUser(req)
}