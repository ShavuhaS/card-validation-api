package card

import (
	"strings"
	"testing"
	"time"
)

func TestExpirationDateValidator(t *testing.T) {
	originalTimeNow := timeNow
	timeNow = func() time.Time {
		return time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	}
	defer func() { timeNow = originalTimeNow }()

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
			name:          "Card number is too short",
			request:       &ValidationRequest{Number: "54321"},
			expectedError: ErrInvalidNumberLength,
		},
		{
			name:          "Card number is too long",
			request:       &ValidationRequest{Number: "1111111111111111111111111111111111"},
			expectedError: ErrInvalidNumberLength,
		},
		{
			name:          "Valid length short card number",
			request:       &ValidationRequest{Number: strings.Repeat("1", MinNumberLength)},
			expectedError: nil,
		},
		{
			name:          "Valid length long card number",
			request:       &ValidationRequest{Number: strings.Repeat("1", MaxNumberLength)},
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

func TestNumberCharsetValidator(t *testing.T) {
	tests := []struct {
		name          string
		request       *ValidationRequest
		expectedError *ValidationError
	}{
		{
			name:          "Card number contains letters",
			request:       &ValidationRequest{Number: "324324233322f"},
			expectedError: ErrInvalidNumberCharacters,
		},
		{
			name:          "Card number contains non-digit symbols",
			request:       &ValidationRequest{Number: "32432843%34_."},
			expectedError: ErrInvalidNumberCharacters,
		},
		{
			name:          "Card number contains non-arabic numerals",
			request:       &ValidationRequest{Number: "ⅤⅢⅤⅤⅤⅤⅤⅢⅤⅤⅤⅤⅤ"},
			expectedError: ErrInvalidNumberCharacters,
		},
		{
			name:          "Card number contains only numeric digits",
			request:       &ValidationRequest{Number: "9876543210"},
			expectedError: nil,
		},
	}

	v := &NumberCharsetValidator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.request)
			if err != tt.expectedError {
				t.Errorf("NumberCharsetValidator.Validate(): expected %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestNumberLuhnValidator(t *testing.T) {
	tests := []struct {
		name          string
		request       *ValidationRequest
		expectedError *ValidationError
	}{
		{
			name:          "Card number doesn't pass the Luhn check",
			request:       &ValidationRequest{Number: "4111111111111112"},
			expectedError: ErrLuhnCheckFailed,
		},
		{
			name:          "Card number doesn't pass the Luhn check",
			request:       &ValidationRequest{Number: "5555555555554445"},
			expectedError: ErrLuhnCheckFailed,
		},
		{
			name:          "Card number doesn't pass the Luhn check",
			request:       &ValidationRequest{Number: "378282246310006"},
			expectedError: ErrLuhnCheckFailed,
		},
		{
			name:          "Card number doesn't pass the Luhn check",
			request:       &ValidationRequest{Number: "6011111111111118"},
			expectedError: ErrLuhnCheckFailed,
		},
		{
			name:          "Card number passes the Luhn check",
			request:       &ValidationRequest{Number: "4242424242424242"},
			expectedError: nil,
		},
		{
			name:          "Card number passes the Luhn check",
			request:       &ValidationRequest{Number: "4111111111111111"},
			expectedError: nil,
		},
		{
			name:          "Card number passes the Luhn check",
			request:       &ValidationRequest{Number: "5555555555554444"},
			expectedError: nil,
		},
		{
			name:          "Card number passes the Luhn check",
			request:       &ValidationRequest{Number: "378282246310005"},
			expectedError: nil,
		},
	}

	v := &NumberLuhnValidator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.request)
			if err != tt.expectedError {
				t.Errorf("NumberLuhnValidator.Validate(): expected %v, got %v", tt.expectedError, err)
			}
		})
	}
}

type mockNetworkStorage struct {
	findIssuer       func(iin string) (CardIssuer, bool)
	getNumberLengths func(issuer CardIssuer) []int
}

func (m *mockNetworkStorage) FindIssuer(iin string) (CardIssuer, bool) {
	return m.findIssuer(iin)
}

func (m *mockNetworkStorage) GetNumberLengths(issuer CardIssuer) []int {
	return m.getNumberLengths(issuer)
}

func TestNumberIINValidator(t *testing.T) {
	mockStorage := &mockNetworkStorage{
		findIssuer: func(iin string) (CardIssuer, bool) {
			if iin == "41111111" {
				return Visa, true
			}
			return Visa, false
		},
	}
	v := NewNumberIINValidator(mockStorage)

	tests := []struct {
		name          string
		request       *ValidationRequest
		expectedError *ValidationError
	}{
		{
			name:          "Valid IIN",
			request:       &ValidationRequest{Number: "4111111111111111"},
			expectedError: nil,
		},
		{
			name:          "Invalid IIN",
			request:       &ValidationRequest{Number: "0000000000000000"},
			expectedError: ErrIINIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.request)
			if err != tt.expectedError {
				t.Errorf("NumberIINValidator.Validate(): expected %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestIINToNumberLengthValidator(t *testing.T) {
	mockStorage := &mockNetworkStorage{
		findIssuer: func(iin string) (CardIssuer, bool) {
			if iin == "41111111" {
				return Visa, true
			}
			return Visa, false
		},
		getNumberLengths: func(issuer CardIssuer) []int {
			if issuer == Visa {
				return []int{16}
			}
			return []int{}
		},
	}
	v := NewIINToNumberLengthValidator(mockStorage)

	tests := []struct {
		name          string
		request       *ValidationRequest
		expectedError *ValidationError
	}{
		{
			name:          "Valid IIN and valid length",
			request:       &ValidationRequest{Number: "4111111111111111"},
			expectedError: nil,
		},
		{
			name:          "Valid IIN but invalid length",
			request:       &ValidationRequest{Number: "411111111111111"},
			expectedError: ErrInvalidNumberLengthForIIN,
		},
		{
			name:          "Invalid IIN",
			request:       &ValidationRequest{Number: "0000000000000000"},
			expectedError: ErrIINIsInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.request)
			if err != tt.expectedError {
				t.Errorf("IINToNumberLengthValidator.Validate(): expected %v, got %v", tt.expectedError, err)
			}
		})
	}
}
