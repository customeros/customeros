package entity

import (
	"time"
)

type CustomFieldTemplateEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Type      string
	Order     int64
	Mandatory bool
	Length    *int64
	Min       *int64
	Max       *int64
}

type CustomFieldTemplateEntities []CustomFieldTemplateEntity
