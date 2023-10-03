package entity

import (
	"time"
)

type PhoneNumberEntity struct {
	Id             string
	E164           string
	RawPhoneNumber string
	Source         DataSource
	SourceOfTruth  DataSource
	AppSource      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PhoneNumberEntities []PhoneNumberEntity
