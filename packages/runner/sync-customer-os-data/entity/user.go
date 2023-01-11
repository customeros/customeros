package entity

import "time"

type UserData struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time

	ExternalId     string
	ExternalSystem string
}
