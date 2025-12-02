package utils

import (
	"booking-service/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var (
	authServiceURL   string
	cinemaServiceURL string
)

func init() {
	authServiceURL = os.Getenv("AUTH_SERVICE_URL")
	cinemaServiceURL = os.Getenv("CINEMA_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:3001"
	}
	if cinemaServiceURL == "" {
		cinemaServiceURL = "http://localhost:3002"
	}
}

func ReserveSeats(seatIDs []uint) error {
	reqBody := map[string][]uint{"seatIds": seatIDs}
	jsonData, _ := json.Marshal(reqBody)
	
	resp, err := http.Post(cinemaServiceURL+"/api/cinema/seats/reserve", "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("failed to reserve seats")
	}
	resp.Body.Close()
	return nil
}

func ReleaseSeats(seatIDs []uint) {
	reqBody := map[string][]uint{"seatIds": seatIDs}
	jsonData, _ := json.Marshal(reqBody)
	http.Post(cinemaServiceURL+"/api/cinema/seats/release", "application/json", bytes.NewBuffer(jsonData))
}

func GetStudioDetails(studioID uint) (*models.StudioDetails, error) {
	url := fmt.Sprintf("%s/api/cinema/studios/%d", cinemaServiceURL, studioID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cinema service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get studio details")
	}

	var studio models.StudioDetails
	if err := json.NewDecoder(resp.Body).Decode(&studio); err != nil {
		return nil, fmt.Errorf("failed to decode studio details")
	}

	return &studio, nil
}

func GetSeatsDetails(seatIDs []int64) ([]models.SeatDetails, error) {
	url := fmt.Sprintf("%s/api/cinema/seats/details", cinemaServiceURL)

	// Convert []int64 to []uint for the request
	uintSeatIDs := make([]uint, len(seatIDs))
	for i, id := range seatIDs {
		uintSeatIDs[i] = uint(id)
	}

	body, _ := json.Marshal(map[string]interface{}{
		"seat_ids": uintSeatIDs,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cinema service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse map[string]string
		json.NewDecoder(resp.Body).Decode(&errResponse)
		return nil, fmt.Errorf("failed to get seat details: %s", errResponse["error"])
	}

	var seats []models.SeatDetails
	if err := json.NewDecoder(resp.Body).Decode(&seats); err != nil {
		return nil, fmt.Errorf("failed to decode seat details")
	}

	return seats, nil
}
