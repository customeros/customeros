package entity

import "time"

type UserData struct {
	Id              string    `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	FirstName       string    `json:"firstName,omitempty"`
	LastName        string    `json:"lastName,omitempty"`
	Email           string    `json:"email,omitempty"`
	PhoneNumber     string    `json:"phoneNumber,omitempty"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"`
	ExternalId      string    `json:"externalId,omitempty"`
	ExternalOwnerId string    `json:"externalOwnerId,omitempty"`

	ExternalSystem string `json:"externalSystem,omitempty"`
	ExternalSyncId string `json:"externalSyncId,omitempty"`
}

func (u UserData) HasPhoneNumber() bool {
	return len(u.PhoneNumber) > 0
}

func (u UserData) HasEmail() bool {
	return len(u.Email) > 0
}
