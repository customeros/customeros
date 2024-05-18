package entity

import "time"

type CountryEntity struct {
	DataLoaderKey

	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	CodeA2    string
	CodeA3    string
	PhoneCode string
}

type CountryEntities []CountryEntity
