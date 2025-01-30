package neo4j_entity

import (
	"time"
)

type TenantProperty string

const (
	TenantPropertyName                TenantProperty = "name"
	TenantPropertyOnboardingCheckedAt TenantProperty = "techOnboardingCheckedAt"
)

type TenantEntity struct {
	Id                   string
	Name                 string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Source               DataSource
	AppSource            string
	Settings             TenantSettingsEntity
	Active               bool
	TenantInternalFields TenantInternalFields
}

type TenantInternalFields struct {
	OnboardingCheckedAt *time.Time
}
