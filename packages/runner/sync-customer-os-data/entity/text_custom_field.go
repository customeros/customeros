package entity

import "time"

type TextCustomField struct {
	Name           string
	Value          string
	ExternalSystem string
	CreatedAt      time.Time
}
