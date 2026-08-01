package card

type ValidationService interface {
	Validate(req *ValidationRequest) *ValidationError
}

type serviceImpl struct {
	validators []CardValidator
}

func NewValidationService() ValidationService {
	return &serviceImpl{
		validators: []CardValidator{
			ExpirationDateValidator{},
			NumberLengthValidator{},
			NumberCharsetValidator{},
			NumberLuhnValidator{},
		},
	}
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
