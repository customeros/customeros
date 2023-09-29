package models

import (
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type UserDataFields struct {
	Name            string
	FirstName       string
	LastName        string
	Internal        bool
	ProfilePhotoUrl string
	Timezone        string
}

type UserFields struct {
	ID             string
	Tenant         string
	UserDataFields UserDataFields
	Source         commonModels.Source
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}
