package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"auth-service/database"
	"auth-service/handlers"
)

func main() {
	database.Init()
	
	r := gin.Default()
	
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/verify", handlers.Verify)
		auth.GET("/google", handlers.GoogleLogin)
		auth.GET("/google/callback", handlers.GoogleCallback)
	}
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "auth-service"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Auth service running on port %s", port)
	r.Run(":" + port)
}