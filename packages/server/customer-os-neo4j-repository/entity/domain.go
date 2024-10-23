package entity

import (
	"time"
)

type DomainProperty string

const (
	DomainPropertyCreatedAt                     DomainProperty = "createdAt"
	DomainPropertyUpdatedAt                     DomainProperty = "updatedAt"
	DomainPropertyDomain                        DomainProperty = "domain"
	DomainPropertySource                        DomainProperty = "source"
	DomainPropertyAppSource                     DomainProperty = "appSource"
	DomainPropertyIsPrimary                     DomainProperty = "primary"
	DomainPropertyPrimaryDomain                 DomainProperty = "primaryDomain"
	DomainPropertyPrimaryDomainCheckRequestedAt DomainProperty = "techPrimaryDomainCheckRequestedAt"
)

type DomainEntity struct {
	DataLoaderKey
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Domain         string
	Source         DataSource
	AppSource      string
	IsPrimary      *bool
	PrimaryDomain  string
	InternalFields DomainInternalFields
}

type DomainInternalFields struct {
	PrimaryDomainCheckRequestedAt *time.Time
}

type DomainEntities []DomainEntity
