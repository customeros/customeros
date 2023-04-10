package models

import (
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type UserDto struct {
	ID        string
	Tenant    string
	FirstName string
	LastName  string
	Name      string
	Source    commonModels.Source
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
