package models

import (
	"time"
)

// StudioDetails contains details about a studio.
type StudioDetails struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	TotalSeats int    `json:"total_seats"`
}

// SeatDetails contains details about a seat.
type SeatDetails struct {
	ID         uint   `json:"id"`
	SeatNumber string `json:"seat_number"`
}

// BookingWithDetails contains a booking and its related studio and seat details.
type BookingWithDetails struct {
	ID          uint           `json:"id"`
	BookingCode string         `json:"booking_code"`
	UserID      *uint          `json:"user_id"`
	UserName    string         `json:"user_name"`
	UserEmail   string         `json:"user_email"`
	QRCode      string         `json:"qr_code"`
	BookingType string         `json:"booking_type"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Studio      *StudioDetails `json:"studio"`
	Seats       []SeatDetails  `json:"seats"`
}
