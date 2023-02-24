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
	Country      string   `json:"country"`
	City         string   `json:"city"`
	Zip          string   `json:"zip"`
	AddressLine1 string   `json:"addressLine1"`
	AddressLine2 string   `json:"addressLine2"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

func MapValidationAddressResponse(validatedAddress *SmartyAddressResponse, error *string, valid bool) ValidationAddressResponse {
	return ValidationAddressResponse{
		Address: MapAddress(validatedAddress), //TODO multiple addresses
		Valid:   valid,
		Error:   error,
	}
}

func MapAddress(address *SmartyAddressResponse) *Address {
	if address == nil {
		return nil
	}
	return &Address{
		//Country:      address.Result.Addresses[0].ApiOutput[0].Components. ,
		Country:      "TODO",
		City:         address.Result.Addresses[0].ApiOutput[0].Components.CityName,
		Zip:          address.Result.Addresses[0].ApiOutput[0].Components.Zipcode,
		AddressLine1: address.Result.Addresses[0].ApiOutput[0].DeliveryLine1,
		//AddressLine2: address.Result.Addresses[0].ApiOutput[0].DeliveryLine1,
		AddressLine2: "TODO",
		Latitude:     &address.Result.Addresses[0].ApiOutput[0].Metadata.Latitude,
		Longitude:    &address.Result.Addresses[0].ApiOutput[0].Metadata.Longitude,
	}
}
