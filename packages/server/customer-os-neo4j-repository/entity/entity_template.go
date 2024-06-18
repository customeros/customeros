package entity

import (
	"time"
)

type EntityTemplateEntity struct {
	Id           string
	Name         string
	Extends      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CustomFields []*CustomFieldTemplateEntity
}

type EntityTemplateEntities []EntityTemplateEntity
