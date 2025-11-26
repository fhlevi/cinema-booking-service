package models

import (
	"time"
	"gorm.io/gorm"
)

type Studio struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"not null"`
	TotalSeats int            `json:"total_seats" gorm:"default:20"`
	Seats      []Seat         `json:"seats,omitempty" gorm:"foreignKey:StudioID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type Seat struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	StudioID    uint           `json:"studio_id" gorm:"not null"`
	SeatNumber  string         `json:"seat_number" gorm:"not null"`
	IsAvailable bool           `json:"is_available" gorm:"default:true"`
	Studio      Studio         `json:"studio,omitempty" gorm:"foreignKey:StudioID"`
	StudioName  string         `json:"studio_name,omitempty" gorm:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type SeatReservationRequest struct {
	SeatIDs []uint `json:"seatIds"`
}