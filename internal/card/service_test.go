package card

import (
	"testing"
	"time"
)

func TestValidationService_ValidatorOrder(t *testing.T) {
	originalTimeNow := timeNow
	timeNow = func() time.Time {
		return time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	}
	defer func() { timeNow = originalTimeNow }()

	service, err := NewValidationService()
	if err != nil {
		t.Fatalf("Failed to initialize ValidationService: %v", err)
	}

	tests := []struct {
		name          string
		request       *ValidationRequest
		expectedError *ValidationError
	}{
		{
			name: "ExpirationDate runs before NumberLength",
			request: &ValidationRequest{
				ExpMonth: 13,
				ExpYear:  2028,
				Number:   "123",
			},
			expectedError: ErrInvalidMonth,
		},
		{
			name: "NumberLength runs before NumberCharset",
			request: &ValidationRequest{
				ExpMonth: 12,
				ExpYear:  2028,
				Number:   "abc",
			},
			expectedError: ErrInvalidNumberLength,
		},
		{
			name: "NumberCharset runs before NumberLuhn",
			request: &ValidationRequest{
				ExpMonth: 12,
				ExpYear:  2028,
				Number:   "411111111111abcd",
			},
			expectedError: ErrInvalidNumberCharacters,
		},
		{
			name: "NumberLuhn runs before NumberIINValidator",
			request: &ValidationRequest{
				ExpMonth: 12,
				ExpYear:  2028,
				Number:   "0000000000000001",
			},
			expectedError: ErrLuhnCheckFailed,
		},
		{
			name: "NumberIINValidator runs before IINToNumberLength",
			request: &ValidationRequest{
				ExpMonth: 12,
				ExpYear:  2028,
				Number:   "0000000000000000",
			},
			expectedError: ErrIINIsInvalid,
		},
		{
			name: "IINToNumberLength runs last",
			request: &ValidationRequest{
				ExpMonth: 12,
				ExpYear:  2028,
				Number:   "51000000000008",
			},
			expectedError: ErrInvalidNumberLengthForIIN,
		},
		{
			name: "All validators pass",
			request: &ValidationRequest{
				ExpMonth: 12,
				ExpYear:  2028,
				Number:   "4111111111111111",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Validate(tt.request)
			if err != tt.expectedError {
				t.Errorf("ValidationService.Validate() returned %v, expected %v", err, tt.expectedError)
			}
		})
	}
}
