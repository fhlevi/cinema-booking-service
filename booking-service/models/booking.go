package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/lib/pq"
)

type Booking struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	BookingCode string         `json:"booking_code" gorm:"uniqueIndex;not null"`
	UserID      *uint          `json:"user_id"`
	UserName    string         `json:"user_name" gorm:"not null"`
	UserEmail   string         `json:"user_email" gorm:"not null"`
	StudioID    uint           `json:"studio_id" gorm:"not null"`
	SeatIDs     pq.Int64Array  `json:"seat_ids" gorm:"type:integer[]"`
	QRCode      string         `json:"qr_code" gorm:"type:text"`
	BookingType string         `json:"booking_type" gorm:"default:online"`
	Status      string         `json:"status" gorm:"default:active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type OnlineBookingRequest struct {
	StudioID uint   `json:"studioId"`
	SeatIDs  []uint `json:"seatIds"`
}

type OfflineBookingRequest struct {
	StudioID      uint   `json:"studioId"`
	SeatIDs       []uint `json:"seatIds"`
	CustomerName  string `json:"customerName"`
	CustomerEmail string `json:"customerEmail"`
}

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type ValidateQRRequest struct {
	BookingCode string `json:"bookingCode"`
}