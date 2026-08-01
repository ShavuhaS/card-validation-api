package card

import (
	"time"
)

const MinNumberLength = 13
const MaxNumberLength = 19

const RealisticYearsUntilExpiration = 20

var timeNow = time.Now

type CardValidator interface {
	Validate(req *ValidationRequest) *ValidationError
}

type ExpirationDateValidator struct{}

func (v ExpirationDateValidator) Validate(req *ValidationRequest) *ValidationError {
	month := req.ExpMonth
	year := req.ExpYear

	if month < 1 || month > 12 {
		return ErrInvalidMonth
	}

	today := timeNow()

	if year < 1000 || year > today.Year()+RealisticYearsUntilExpiration {
		return ErrInvalidYear
	}

	if year < today.Year() || (year == today.Year() && month < int(today.Month())) {
		return ErrCardExpired
	}

	return nil
}

type NumberLengthValidator struct{}

func (v NumberLengthValidator) Validate(req *ValidationRequest) *ValidationError {
	length := len(req.Number)
	if length < MinNumberLength || length > MaxNumberLength {
		return ErrInvalidNumberLength
	}
	return nil
}

type NumberCharsetValidator struct{}

func (v NumberCharsetValidator) Validate(req *ValidationRequest) *ValidationError {
	for _, ch := range req.Number {
		if ch < '0' || ch > '9' {
			return ErrInvalidNumberCharacters
		}
	}
	return nil
}

type NumberLuhnValidator struct{}

func (v NumberLuhnValidator) Validate(req *ValidationRequest) *ValidationError {
	checksum := 0
	parity := len(req.Number) % 2
	for i, ch := range req.Number {
		digit := int(ch - '0')
		if (i % 2) == parity {
			digit *= 2
			if digit >= 10 {
				digit -= 9
			}
		}
		checksum += digit
	}

	if checksum%10 != 0 {
		return ErrLuhnCheckFailed
	}

	return nil
}
