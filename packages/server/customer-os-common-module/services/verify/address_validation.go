package verify

import (
	"context"

	"github.com/sirupsen/logrus"
	international_street "github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	extract "github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	street "github.com/smartystreets/smartystreets-go-sdk/us-street-api"
)

type addressValidationService struct {
	USClient   *extract.Client
	IntlClient *international_street.Client
}

func (s *verifyService) ValidateUsAddress(address string) (*extract.Lookup, error) {
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

func (s *verifyService) ValidateInternationalAddress(address, country string) (*international_street.Lookup, error) {
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
