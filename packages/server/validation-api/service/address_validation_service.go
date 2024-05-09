package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/sirupsen/logrus"
	international_street "github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	extract "github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	street "github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

type AddressValidationService interface {
	ValidateUsAddress(address string) (*extract.Lookup, error)
	ValidateInternationalAddress(address, country string) (*international_street.Lookup, error)
}

type addressValidationService struct {
	Services   *Services
	USClient   *extract.Client
	IntlClient *international_street.Client
}

func NewAddressValidationService(config *config.Config, services *Services) AddressValidationService {
	return &addressValidationService{
		Services:   services,
		USClient:   wireup.BuildUSExtractAPIClient(wireup.SecretKeyCredential(config.SmartyConfig.AuthId, config.SmartyConfig.AuthToken)),
		IntlClient: wireup.BuildInternationalStreetAPIClient(wireup.SecretKeyCredential(config.SmartyConfig.AuthId, config.SmartyConfig.AuthToken)),
	}
}

func (s *addressValidationService) ValidateUsAddress(address string) (*extract.Lookup, error) {
	lookup := &extract.Lookup{
		Text:                    address,
		Aggressive:              true,
		AddressesWithLineBreaks: false,
		AddressesPerLine:        1,
		MatchStrategy:           street.MatchEnhanced,
	}

	if err := s.USClient.SendLookupWithContext(context.Background(), lookup); err != nil {
		logrus.Errorf("Error sending batch: {%v}", err)
		return nil, err
	}

	return lookup, nil
}

func (s *addressValidationService) ValidateInternationalAddress(address, country string) (*international_street.Lookup, error) {
	lookup := &international_street.Lookup{
		Freeform: address,
		Country:  country,
	}

	if err := s.IntlClient.SendLookupWithContext(context.Background(), lookup); err != nil {
		logrus.Errorf("Error sending batch: {%v}", err.Error())
		return nil, err
	}

	return lookup, nil
}
