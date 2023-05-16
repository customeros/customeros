package dto

type ValidationAddressRequest struct {
	Address string `json:"address"`
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

func MapValidationAddressResponse(validatedAddress *SmartyAddressResponse, error *string, valid bool) ValidationAddressResponse {
	return ValidationAddressResponse{
		Address: MapAddress(validatedAddress), //TODO handle multiple addresses
		Valid:   valid,
		Error:   error,
	}
}

func MapAddress(address *SmartyAddressResponse) *Address {
	if address == nil {
		return nil
	}
	verifiedAddress := address.Result.Addresses[0].ApiOutput[0]
	for _, v := range address.Result.Addresses {
		if v.Verified {
			verifiedAddress = v.ApiOutput[0]
			break
		}
	}
	returnedAddress := Address{
		Region:       verifiedAddress.Components.StateAbbreviation,
		District:     verifiedAddress.Metadata.CountyName,
		Locality:     verifiedAddress.Components.CityName,
		HouseNumber:  verifiedAddress.Components.PrimaryNumber,
		Street:       verifiedAddress.Components.StreetName,
		Zip:          verifiedAddress.Components.Zipcode,
		PlusFour:     verifiedAddress.Components.Plus4Code,
		AddressLine1: verifiedAddress.DeliveryLine1,
		Latitude:     verifiedAddress.Metadata.Latitude,
		Longitude:    verifiedAddress.Metadata.Longitude,
		TimeZone:     verifiedAddress.Metadata.TimeZone,
		UtcOffset:    verifiedAddress.Metadata.UtcOffset,
	}

	if returnedAddress.Street != "" && verifiedAddress.Components.StreetSuffix != "" {
		returnedAddress.Street += " " + verifiedAddress.Components.StreetSuffix
	}

	return &returnedAddress
}
