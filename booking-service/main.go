package main

import (
	"log"
	"os"

	"booking-service/database"
	"booking-service/handlers"
	"github.com/gin-gonic/gin"
)



func main() {
	database.Init()
	
	r := gin.Default()
	
	booking := r.Group("/api/booking")
	{
		booking.POST("/online", handlers.AuthMiddleware(), handlers.CreateOnlineBooking)
		booking.POST("/offline", handlers.CreateOfflineBooking)
		booking.POST("/validate", handlers.ValidateQRCode)
		booking.GET("/my-bookings", handlers.AuthMiddleware(), handlers.GetUserBookings)
	}
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "booking-service"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Booking service running on port %s", port)
	r.Run(":" + port)
}