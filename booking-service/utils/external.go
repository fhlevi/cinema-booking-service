package utils

import (
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