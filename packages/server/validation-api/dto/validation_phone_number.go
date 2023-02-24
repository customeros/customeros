package dto

type ValidationPhoneNumberRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Country     string `json:"country"`
}

type ValidationPhoneNumberResponse struct {
	E164  *string `json:"e164"`
	Valid bool    `json:"valid"`
	Error *string `json:"error"`
}

func MapValidationPhoneNumberResponse(e164 *string, error *string, valid bool) ValidationPhoneNumberResponse {
	return ValidationPhoneNumberResponse{
		E164:  e164,
		Valid: valid,
		Error: error,
	}
}
