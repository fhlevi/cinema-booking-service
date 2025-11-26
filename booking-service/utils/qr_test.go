package utils

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQRCode(t *testing.T) {
	tests := []struct {
		name         string
		bookingCode  string
		studioID     uint
		seatIDs      []uint
		userID       *uint
		customerName string
		expectError  bool
	}{
		{
			name:        "Valid QR code with user ID",
			bookingCode: "BOOK123",
			studioID:    1,
			seatIDs:     []uint{1, 2, 3},
			userID:      func() *uint { id := uint(1); return &id }(),
			expectError: false,
		},
		{
			name:         "Valid QR code with customer name",
			bookingCode:  "BOOK456",
			studioID:     2,
			seatIDs:      []uint{4, 5},
			customerName: "John Doe",
			expectError:  false,
		},
		{
			name:        "Empty booking code",
			bookingCode: "",
			studioID:    1,
			seatIDs:     []uint{1},
			userID:      func() *uint { id := uint(1); return &id }(),
			expectError: false, // QR generation should still work
		},
		{
			name:        "Empty seat IDs",
			bookingCode: "BOOK789",
			studioID:    1,
			seatIDs:     []uint{},
			userID:      func() *uint { id := uint(1); return &id }(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qrCode, err := GenerateQRCode(tt.bookingCode, tt.studioID, tt.seatIDs, tt.userID, tt.customerName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, qrCode)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, qrCode)

				// Verify QR code format
				assert.True(t, strings.HasPrefix(qrCode, "data:image/png;base64,"))

				// Extract and decode base64 data
				base64Data := strings.TrimPrefix(qrCode, "data:image/png;base64,")
				_, err := base64.StdEncoding.DecodeString(base64Data)
				assert.NoError(t, err, "QR code should contain valid base64 data")
			}
		})
	}
}

func TestQRCodeDataStructure(t *testing.T) {
	bookingCode := "TEST123"
	studioID := uint(1)
	seatIDs := []uint{1, 2, 3}
	userID := uint(1)
	customerName := "Test User"

	// Test with user ID
	qrCode, err := GenerateQRCode(bookingCode, studioID, seatIDs, &userID, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, qrCode)

	// Test with customer name
	qrCode2, err := GenerateQRCode(bookingCode, studioID, seatIDs, nil, customerName)
	assert.NoError(t, err)
	assert.NotEmpty(t, qrCode2)

	// Verify different QR codes are generated for different inputs
	assert.NotEqual(t, qrCode, qrCode2)
}

func TestQRCodeConsistency(t *testing.T) {
	bookingCode := "CONSISTENT123"
	studioID := uint(1)
	seatIDs := []uint{1, 2, 3}
	userID := uint(1)

	// Generate QR code twice with same inputs
	qrCode1, err1 := GenerateQRCode(bookingCode, studioID, seatIDs, &userID, "")
	qrCode2, err2 := GenerateQRCode(bookingCode, studioID, seatIDs, &userID, "")

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// QR codes should be different due to timestamp
	assert.NotEqual(t, qrCode1, qrCode2)
}
