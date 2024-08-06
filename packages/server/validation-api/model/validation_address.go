package model

import (
	international_street "github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	extract "github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
)

type ValidationAddressRequest struct {
	Address       string `json:"address"`
	Country       string `json:"country"`
	International bool   `json:"international"`
}
type ValidationAddressResponse struct {
	Address *Address `json:"address"`
	Valid   bool     `json:"valid"`
	Error   *string  `json:"error"`
}
type Address struct {
	Country      string  `json:"country"`
	Region       string  `json:"region"`
	District     string  `json:"district"`
	Locality     string  `json:"locality"`
	Street       string  `json:"street"`
	Zip          string  `json:"zip"`
	PostalCode   string  `json:"postalCode"`
	AddressLine1 string  `json:"addressLine1"`
	AddressLine2 string  `json:"addressLine2"`
	AddressType  string  `json:"addressType"`
	HouseNumber  string  `json:"houseNumber"`
	PlusFour     string  `json:"plusFour"`
	Commercial   bool    `json:"commercial"`
	Predirection string  `json:"predirection"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	TimeZone     string  `json:"timeZone"`
	UtcOffset    int     `json:"utcOffset"`
}

func MapValidationNoAddressResponse(error *string) ValidationAddressResponse {
	return ValidationAddressResponse{
		Address: nil,
		Valid:   false,
		Error:   error,
	}
}

func MapValidationUsAddressResponse(lookup *extract.Lookup, error *string, valid bool) ValidationAddressResponse {
	return ValidationAddressResponse{
		Address: MapUSAddress(lookup),
		Valid:   valid,
		Error:   error,
	}
}

func MapValidationInternationalAddressResponse(lookup *international_street.Lookup, error *string, valid bool) ValidationAddressResponse {
	return ValidationAddressResponse{
		Address: MapInternationalAddress(lookup),
		Valid:   valid,
		Error:   error,
	}
}

func MapUSAddress(lookup *extract.Lookup) *Address {
	if lookup == nil {
		return nil
	}
	verifiedAddress := lookup.Result.Addresses[0].APIOutput[0]
	for _, v := range lookup.Result.Addresses {
		if v.Verified {
			verifiedAddress = v.APIOutput[0]
			break
		}
	}
	returnedAddress := Address{
		Region:       verifiedAddress.Components.StateAbbreviation,
		District:     verifiedAddress.Metadata.CountyName,
		Locality:     verifiedAddress.Components.CityName,
		HouseNumber:  verifiedAddress.Components.PrimaryNumber,
		Street:       verifiedAddress.Components.StreetName,
		Zip:          verifiedAddress.Components.ZIPCode,
		PlusFour:     verifiedAddress.Components.Plus4Code,
		AddressLine1: verifiedAddress.DeliveryLine1,
		Latitude:     verifiedAddress.Metadata.Latitude,
		Longitude:    verifiedAddress.Metadata.Longitude,
		TimeZone:     verifiedAddress.Metadata.TimeZone,
		UtcOffset:    int(verifiedAddress.Metadata.UTCOffset),
	}

	if returnedAddress.Street != "" && verifiedAddress.Components.StreetSuffix != "" {
		returnedAddress.Street += " " + verifiedAddress.Components.StreetSuffix
	}

	return &returnedAddress
}

func MapInternationalAddress(lookup *international_street.Lookup) *Address {
	if lookup == nil || len(lookup.Results) == 0 {
		return nil
	}
	verifiedAddress := lookup.Results[0]
	for _, v := range lookup.Results {
		if v.Analysis.VerificationStatus == "Verified" {
			verifiedAddress = v
			break
		}
	}
	returnedAddress := Address{
		Country:     verifiedAddress.Components.CountryISO3,
		Region:      verifiedAddress.Components.AdministrativeArea,
		District:    verifiedAddress.Components.SubAdministrativeArea,
		Locality:    verifiedAddress.Components.Locality,
		HouseNumber: verifiedAddress.Components.Premise,
		Street:      verifiedAddress.Components.Thoroughfare,
		PostalCode:  verifiedAddress.Components.PostalCode,
		Latitude:    verifiedAddress.Metadata.Latitude,
		Longitude:   verifiedAddress.Metadata.Longitude,
	}

	return &returnedAddress
}
