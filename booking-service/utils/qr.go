package utils

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(bookingCode string, studioID uint, seatIDs []uint, userID *uint, customerName string) (string, error) {
	qrData := map[string]interface{}{
		"bookingCode": bookingCode,
		"studioId":    studioID,
		"seatIds":     seatIDs,
		"timestamp":   time.Now().Format(time.RFC3339),
	}
	
	if userID != nil {
		qrData["userId"] = *userID
	} else {
		qrData["customerName"] = customerName
	}

	qrDataJSON, err := json.Marshal(qrData)
	if err != nil {
		return "", err
	}

	qrCodeBytes, err := qrcode.Encode(string(qrDataJSON), qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrCodeBytes), nil
}