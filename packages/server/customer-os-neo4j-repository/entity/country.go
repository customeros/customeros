package entity

import "time"

type CountryEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Name      string
	CodeA2    string
	CodeA3    string
	PhoneCode string

	DataloaderKey string
}

type CountryEntities []CountryEntity
