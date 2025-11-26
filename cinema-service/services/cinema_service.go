package services

import (
	"fmt"

	"cinema-service/database"
	"cinema-service/models"
)

func GetAllStudios() ([]models.Studio, error) {
	var studios []models.Studio
	result := database.DB.Find(&studios)
	if result.Error != nil {
		return nil, result.Error
	}
	return studios, nil
}

func GetStudioSeats(studioID string) ([]models.Seat, error) {
	var seats []models.Seat
	result := database.DB.Preload("Studio").Where("studio_id = ?", studioID).Order("seat_number").Find(&seats)
	if result.Error != nil {
		return nil, result.Error
	}

	// Set studio name for response
	for i := range seats {
		seats[i].StudioName = seats[i].Studio.Name
	}

	return seats, nil
}

func ReserveSeats(seatIDs []uint) error {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if seats are available
	var availableSeats []models.Seat
	result := tx.Where("id IN ? AND is_available = ?", seatIDs, true).Find(&availableSeats)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if len(availableSeats) != len(seatIDs) {
		tx.Rollback()
		return fmt.Errorf("some seats are not available")
	}

	// Reserve seats
	result = tx.Model(&models.Seat{}).Where("id IN ?", seatIDs).Update("is_available", false)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

func ReleaseSeats(seatIDs []uint) error {
	result := database.DB.Model(&models.Seat{}).Where("id IN ?", seatIDs).Update("is_available", true)
	return result.Error
}
