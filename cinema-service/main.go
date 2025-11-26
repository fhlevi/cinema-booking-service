package main

import (
	"log"
	"os"

	"cinema-service/database"
	"cinema-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()

	r := gin.Default()

	cinema := r.Group("/api/cinema")
	{
		cinema.GET("/studios", handlers.GetStudios)
		cinema.GET("/studios/:id/seats", handlers.GetStudioSeats)
		cinema.POST("/seats/reserve", handlers.ReserveSeats)
		cinema.POST("/seats/release", handlers.ReleaseSeats)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "cinema-service"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Cinema service running on port %s", port)
	r.Run(":" + port)
}
