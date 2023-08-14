package models

import (
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type ContactDto struct {
	ID              string
	Tenant          string
	FirstName       string
	LastName        string
	Name            string
	Prefix          string
	Description     string
	Timezone        string
	ProfilePhotoUrl string
	Source          commonModels.Source
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
