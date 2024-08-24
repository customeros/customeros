package entity

import (
	"time"
)

type TenantEntity struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    DataSource
	AppSource string
	Settings  TenantSettingsEntity
	Active    bool
}
