package model

type ValidationPhoneNumberRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Country     string `json:"country"`
}

type ValidationPhoneNumberResponse struct {
	E164      *string `json:"e164"`
	CountryA2 *string `json:"countryA2"`
	Valid     bool    `json:"valid"`
	Error     *string `json:"error"`
}

func MapValidationPhoneNumberResponse(e164 *string, countryA2 *string, error *string, valid bool) ValidationPhoneNumberResponse {
	return ValidationPhoneNumberResponse{
		E164:      e164,
		Valid:     valid,
		Error:     error,
		CountryA2: countryA2,
	}
}
