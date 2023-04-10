package models

import (
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type OrganizationDto struct {
	ID          string
	Tenant      string
	Name        string
	Description string
	Website     string
	Industry    string
	IsPublic    bool
	Source      commonModels.Source
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}
