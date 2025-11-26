package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"booking-service/models"
	"booking-service/services"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		token := strings.Replace(authHeader, "Bearer ", "", 1)
		
		reqBody := map[string]string{"token": token}
		jsonData, _ := json.Marshal(reqBody)
		
		authServiceURL := os.Getenv("AUTH_SERVICE_URL")
		if authServiceURL == "" {
			authServiceURL = "http://localhost:3001"
		}
		
		resp, err := http.Post(authServiceURL+"/api/auth/verify", "application/json", bytes.NewBuffer(jsonData))
		if err != nil || resp.StatusCode != 200 {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		var result struct {
			User  models.User `json:"user"`
			Valid bool `json:"valid"`
		}
		json.NewDecoder(resp.Body).Decode(&result)

		if !result.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user", result.User)
		c.Next()
	}
}

func CreateOnlineBooking(c *gin.Context) {
	var req models.OnlineBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, _ := c.Get("user")
	userObj := user.(models.User)

	booking, err := services.CreateOnlineBooking(req, userObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"booking": booking, "qrCode": booking.QRCode})
}

func CreateOfflineBooking(c *gin.Context) {
	var req models.OfflineBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	booking, err := services.CreateOfflineBooking(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"booking": booking, "qrCode": booking.QRCode})
}

func ValidateQRCode(c *gin.Context) {
	var req models.ValidateQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	booking, err := services.ValidateQRCode(req.BookingCode)
	if err != nil {
		if err.Error() == "invalid or used ticket" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert pq.Int64Array back to []uint for response
	seatIDs := make([]uint, len(booking.SeatIDs))
	for i, id := range booking.SeatIDs {
		seatIDs[i] = uint(id)
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"booking": gin.H{
			"bookingCode":  booking.BookingCode,
			"studioId":     booking.StudioID,
			"seatIds":      seatIDs,
			"customerName": booking.UserName,
			"bookingType":  booking.BookingType,
		},
	})
}

func GetUserBookings(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(models.User)

	bookings, err := services.GetUserBookings(userObj.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}