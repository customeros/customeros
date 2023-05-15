package service

import (
	"context"
	"github.com/nyaruka/phonenumbers"
)

type PhoneNumberValidationService interface {
	ValidatePhoneNumber(ctx context.Context, countryCodeA2 string, phoneNumber string) (*string, *string, error)
}

type phoneNumberValidationService struct {
	Services *Services
}

func NewPhoneNumberValidationService(services *Services) PhoneNumberValidationService {
	return &phoneNumberValidationService{
		Services: services,
	}
}

func (s *phoneNumberValidationService) ValidatePhoneNumber(ctx context.Context, countryCodeA2 string, phoneNumber string) (*string, *string, error) {
	num, err := phonenumbers.Parse(phoneNumber, countryCodeA2)
	if err != nil {
		return nil, nil, err
	}
	if !phonenumbers.IsValidNumber(num) {
		return nil, nil, nil
	} else {
		e164 := phonenumbers.Format(num, phonenumbers.E164)
		extractedCountryCodeA2 := phonenumbers.GetRegionCodeForNumber(num)
		return &e164, &extractedCountryCodeA2, nil
	}
}
