package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/dto"
	extract "github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	street "github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
	"log"
)

type AddressValidationService interface {
	ValidateAddress(address string) (*dto.SmartyAddressResponse, error)
}

type addressValidationService struct {
	Services *Services
	Client   *extract.Client
}

func NewAddressValidationService(config *config.Config, services *Services) AddressValidationService {
	return &addressValidationService{
		Services: services,
		Client:   wireup.BuildUSExtractAPIClient(wireup.SecretKeyCredential(config.Smarty.AuthId, config.Smarty.AuthToken)),
	}
}

func (s *addressValidationService) ValidateAddress(address string) (*dto.SmartyAddressResponse, error) {

	lookup := &extract.Lookup{
		Text:                    address,
		Aggressive:              true,
		AddressesWithLineBreaks: false,
		AddressesPerLine:        1,
		MatchStrategy:           street.MatchEnhanced,
	}

	if err := s.Client.SendLookupWithContext(context.Background(), lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	d := new(dto.SmartyAddressResponse)
	err := json.Unmarshal([]byte(DumpJSON(lookup)), &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func DumpJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}

	var indent bytes.Buffer
	err = json.Indent(&indent, b, "", "  ")
	if err != nil {
		return err.Error()
	}
	return indent.String()
}
