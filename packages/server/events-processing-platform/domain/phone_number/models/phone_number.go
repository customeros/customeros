package models

import "time"

type PhoneNumber struct {
	ID             string    `json:"id"`
	RawPhoneNumber string    `json:"rawPhoneNumber"`
	E164           string    `json:"e164"`
	Source         string    `json:"source"`
	SourceOfTruth  string    `json:"sourceOfTruth"`
	AppSource      string    `json:"appSource"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func NewPhoneNumber() *PhoneNumber {
	return &PhoneNumber{}
}
