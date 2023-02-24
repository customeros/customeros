package service

import (
	"context"
	"github.com/nyaruka/phonenumbers"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
)

type PhoneNumberValidationService interface {
	ValidatePhoneNumber(ctx context.Context, countryCodeA3 string, phoneNumber string) (*string, error)
}

type phoneNumberValidationService struct {
	Services *Services
}

func NewPhoneNumberValidationService(config *config.Config, services *Services) PhoneNumberValidationService {
	return &phoneNumberValidationService{
		Services: services,
	}
}

func (s *phoneNumberValidationService) ValidatePhoneNumber(ctx context.Context, countryCodeA3 string, phoneNumber string) (*string, error) {
	//country, err := s.Services.CommonServices.CountryService.GetCountryByCodeA3(ctx, countryCodeA3)
	//if err != nil {
	//	return nil, err
	//}
	//if country == nil {
	//	return nil, nil
	//}

	num, err := phonenumbers.Parse(phoneNumber, countryCodeA3)
	if err != nil {
		return nil, err
	}
	if !phonenumbers.IsValidNumber(num) {
		return nil, nil
	} else {
		e164 := phonenumbers.Format(num, phonenumbers.E164)
		return &e164, nil
	}
}
