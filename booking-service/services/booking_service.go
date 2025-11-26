package services

import (
	"fmt"

	"booking-service/database"
	"booking-service/models"
	"booking-service/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func CreateOnlineBooking(req models.OnlineBookingRequest, user models.User) (*models.Booking, error) {
	return createBooking(req.StudioID, req.SeatIDs, &user.ID, user.Name, user.Email, "online")
}

func CreateOfflineBooking(req models.OfflineBookingRequest) (*models.Booking, error) {
	return createBooking(req.StudioID, req.SeatIDs, nil, req.CustomerName, req.CustomerEmail, "offline")
}

func createBooking(studioID uint, seatIDs []uint, userID *uint, userName, userEmail, bookingType string) (*models.Booking, error) {
	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Reserve seats in cinema service
	err := utils.ReserveSeats(seatIDs)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	bookingCode := uuid.New().String()
	
	// Generate QR code
	qrCode, err := utils.GenerateQRCode(bookingCode, studioID, seatIDs, userID, userName)
	if err != nil {
		tx.Rollback()
		utils.ReleaseSeats(seatIDs)
		return nil, fmt.Errorf("failed to generate QR code")
	}

	// Convert []uint to pq.Int64Array
	seatIDsInt64 := make(pq.Int64Array, len(seatIDs))
	for i, id := range seatIDs {
		seatIDsInt64[i] = int64(id)
	}

	// Create booking
	booking := models.Booking{
		BookingCode: bookingCode,
		UserID:      userID,
		UserName:    userName,
		UserEmail:   userEmail,
		StudioID:    studioID,
		SeatIDs:     seatIDsInt64,
		QRCode:      qrCode,
		BookingType: bookingType,
		Status:      "active",
	}

	result := tx.Create(&booking)
	if result.Error != nil {
		tx.Rollback()
		utils.ReleaseSeats(seatIDs)
		return nil, fmt.Errorf("failed to create booking")
	}

	if err := tx.Commit().Error; err != nil {
		utils.ReleaseSeats(seatIDs)
		return nil, fmt.Errorf("failed to commit transaction")
	}

	return &booking, nil
}

func ValidateQRCode(bookingCode string) (*models.Booking, error) {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var booking models.Booking
	result := tx.Where("booking_code = ? AND status = ?", bookingCode, "active").First(&booking)
	if result.Error != nil {
		tx.Rollback()
		return nil, fmt.Errorf("invalid or used ticket")
	}

	// Mark as used
	result = tx.Model(&booking).Update("status", "used")
	if result.Error != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update booking status")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction")
	}

	return &booking, nil
}

func GetUserBookings(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	result := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}

	return bookings, nil
}