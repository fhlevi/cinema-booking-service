package database

import (
	"fmt"
	"log"
	"os"

	"cinema-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.Studio{}, &models.Seat{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Add unique constraint
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_studio_seat ON seats(studio_id, seat_number)")

	initializeData()
	log.Println("Cinema database initialized with GORM")
}

func initializeData() {
	var count int64
	DB.Model(&models.Studio{}).Count(&count)

	if count == 0 {
		for i := 1; i <= 5; i++ {
			studio := models.Studio{
				Name:       fmt.Sprintf("Studio %d", i),
				TotalSeats: 20,
			}
			DB.Create(&studio)

			for j := 1; j <= 20; j++ {
				seat := models.Seat{
					StudioID:    studio.ID,
					SeatNumber:  fmt.Sprintf("A%d", j),
					IsAvailable: true,
				}
				DB.Create(&seat)
			}
		}
		log.Println("Initialized 5 studios with 20 seats each")
	}
}