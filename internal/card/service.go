package card

import "fmt"

type ValidationService interface {
	Validate(req *ValidationRequest) *ValidationError
}

type serviceImpl struct {
	validators []CardValidator
}

func NewValidationService() (ValidationService, error) {
	networkStorage, err := NewNetworkStorage()
	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate network storage: %v", err)
	}

	return &serviceImpl{
		validators: []CardValidator{
			ExpirationDateValidator{},
			NumberLengthValidator{},
			NumberCharsetValidator{},
			NumberLuhnValidator{},
			NewNumberIINValidator(networkStorage),
			NewIINToNumberLengthValidator(networkStorage),
		},
	}, nil
}

func (s *serviceImpl) Validate(req *ValidationRequest) *ValidationError {
	for _, v := range s.validators {
		err := v.Validate(req)
		if err != nil {
			return err
		}
	}
	return nil
}
