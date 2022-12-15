package entity

import "time"

type ContactEntity struct {
	Id             string
	Title          string
	FirstName      string
	LastName       string
	Label          string
	Notes          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Readonly       bool
	ExternalId     string
	ExternalSystem string
}
