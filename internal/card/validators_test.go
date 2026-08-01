package card

import (
	"strings"
	"testing"
	"time"
)

func TestExpirationDateValidator(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	}

	tests := []struct {
		name          string
		request       *ValidationRequest
		expectedError *ValidationError
	}{
		{
			name:          "Invalid month (non-positive)",
			request:       &ValidationRequest{ExpMonth: 0, ExpYear: 2026},
			expectedError: ErrInvalidMonth,
		},
		{
			name:          "Invalid month (more than 12)",
			request:       &ValidationRequest{ExpMonth: 13, ExpYear: 2020},
			expectedError: ErrInvalidMonth,
		},
		{
			name:          "Invalid year (less than 4 digits)",
			request:       &ValidationRequest{ExpMonth: 1, ExpYear: 999},
			expectedError: ErrInvalidYear,
		},
		{
			name:          "Invalid year (too far in the future)",
			request:       &ValidationRequest{ExpMonth: 1, ExpYear: 2050},
			expectedError: ErrInvalidYear,
		},
		{
			name:          "Card expired last year",
			request:       &ValidationRequest{ExpMonth: 12, ExpYear: 2025},
			expectedError: ErrCardExpired,
		},
		{
			name:          "Card expired this year",
			request:       &ValidationRequest{ExpMonth: 4, ExpYear: 2026},
			expectedError: ErrCardExpired,
		},
		{
			name:          "Valid date this month",
			request:       &ValidationRequest{ExpMonth: 8, ExpYear: 2026},
			expectedError: nil,
		},
		{
			name:          "Valid date this year",
			request:       &ValidationRequest{ExpMonth: 12, ExpYear: 2026},
			expectedError: nil,
		},
		{
			name:          "Valid date in distant future",
			request:       &ValidationRequest{ExpMonth: 3, ExpYear: 2046},
			expectedError: nil,
		},
	}

	v := ExpirationDateValidator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.request)
			if err != tt.expectedError {
				t.Errorf("ExpirationDateValidator.Validate(): expected %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestNumberLengthValidator(t *testing.T) {
	tests := []struct {
		name          string
		request       *ValidationRequest
		expectedError *ValidationError
	}{
		{
			name: "Card number is too short",
			request: &ValidationRequest{Number: "54321"},
			expectedError: ErrInvalidNumberLength,
		},
		{
			name: "Card number is too long",
			request: &ValidationRequest{Number: "1111111111111111111111111111111111"},
			expectedError: ErrInvalidNumberLength,
		},
		{
			name: "Valid length short card number",
			request: &ValidationRequest{Number: strings.Repeat("1", MinNumberLength)},
			expectedError: nil,
		},
		{
			name: "Valid length long card number",
			request: &ValidationRequest{Number: strings.Repeat("1", MaxNumberLength)},
			expectedError: nil,
		},
	}

	v := &NumberLengthValidator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.request)
			if err != tt.expectedError {
				t.Errorf("NumberLengthValidator.Validate(): expected %v, got %v", tt.expectedError, err)
			}
		})
	}
}

