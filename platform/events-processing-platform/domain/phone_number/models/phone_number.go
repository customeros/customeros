package models

type PhoneNumber struct {
	ID             string `json:"id"`
	RawPhoneNumber string `json:"rawPhoneNumber"`
	E164           string `json:"e164"`
}

func NewPhoneNumber() *PhoneNumber {
	return &PhoneNumber{}
}
