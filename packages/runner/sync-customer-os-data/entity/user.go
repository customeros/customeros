package entity

import "time"

type UserData struct {
	Id          string
	Name        string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	ExternalId      string
	ExternalOwnerId string
	ExternalSystem  string

	ExternalSyncId string
}

func (u UserData) HasPhoneNumber() bool {
	return len(u.PhoneNumber) > 0
}

func (u UserData) HasEmail() bool {
	return len(u.Email) > 0
}
