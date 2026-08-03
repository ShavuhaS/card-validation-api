package card

import "fmt"

type ValidationRequest struct {
	Number   string `json:"cardNumber"`
	ExpMonth int    `json:"expMonth"`
	ExpYear  int    `json:"expYear"`
}

type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("CardValidationError #%v: %v", err.Code, err.Message)
}

type ValidationResponse struct {
	Valid bool             `json:"valid"`
	Error *ValidationError `json:"error,omitempty"`
}

func NewValidationResponse(error *ValidationError) ValidationResponse {
	return ValidationResponse{
		Valid: error == nil,
		Error: error,
	}
}
