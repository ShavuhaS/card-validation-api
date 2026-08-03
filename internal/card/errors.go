package card

import "fmt"

var ErrParsingBody = &ValidationError{
	Code:    "400",
	Message: "Error parsing body",
}

var ErrInvalidMonth = &ValidationError{
	Code:    "001",
	Message: "Invalid expiration month value",
}

var ErrInvalidYear = &ValidationError{
	Code:    "002",
	Message: "Invalid expiration year value",
}

var ErrCardExpired = &ValidationError{
	Code:    "003",
	Message: "Credit card has expired",
}

var ErrInvalidNumberLength = &ValidationError{
	Code:    "004",
	Message: fmt.Sprintf("Invalid card number length. It must be from %v to %v digits long", MinNumberLength, MaxNumberLength),
}

var ErrInvalidNumberCharacters = &ValidationError{
	Code:    "005",
	Message: "Invalid card number characters. Card number may only contain numeric digits 0-9",
}

var ErrMIICheckFailed = &ValidationError{
	Code:    "006",
	Message: "Invalid card number. MII is not supported",
}

var ErrLuhnCheckFailed = &ValidationError{
	Code:    "007",
	Message: "Invalid card number. Luhn's algorithm check failed",
}

var ErrIINIsInvalid = &ValidationError{
	Code:    "008",
	Message: "Invalid card number. IIN is not supported",
}

var ErrInvalidNumberLengthForIIN = &ValidationError{
	Code:    "009",
	Message: "Invalid card number. The length doesn't match the IIN",
}
