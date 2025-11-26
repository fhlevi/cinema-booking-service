package handlers

import (
	"net/http"

	"cinema-service/models"
	"cinema-service/services"

	"github.com/gin-gonic/gin"
)

func GetStudios(c *gin.Context) {
	studios, err := services.GetAllStudios()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch studios"})
		return
	}

	c.JSON(http.StatusOK, studios)
}

func GetStudioSeats(c *gin.Context) {
	studioID := c.Param("id")

	seats, err := services.GetStudioSeats(studioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seats"})
		return
	}

	c.JSON(http.StatusOK, seats)
}

func ReserveSeats(c *gin.Context) {
	var req models.SeatReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := services.ReserveSeats(req.SeatIDs)
	if err != nil {
		if err.Error() == "some seats are not available" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reserve seats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seats reserved successfully"})
}

func ReleaseSeats(c *gin.Context) {
	var req models.SeatReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := services.ReleaseSeats(req.SeatIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to release seats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seats released successfully"})
}
