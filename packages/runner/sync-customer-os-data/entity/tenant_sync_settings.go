package entity

import (
	"database/sql/driver"
	"time"
)

type AirbyteSource string

const (
	HUBSPOT AirbyteSource = "hubspot"
	ZENDESK AirbyteSource = "zendesk"
)

type TenantSyncSettings struct {
	ID        uint          `gorm:"primarykey"`
	CreatedAt time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
	Tenant    string        `gorm:"column:tenant;not null"`
	Source    AirbyteSource `gorm:"type:airbyte_source;column:source;not null"`
	Enabled   bool          `gorm:"column:enabled;not null"`
}

func (as *AirbyteSource) Scan(value interface{}) error {
	*as = AirbyteSource(value.(string))
	return nil
}

func (as AirbyteSource) Value() (driver.Value, error) {
	return string(as), nil
}

type TenantSyncSettingsList []TenantSyncSettings

func (TenantSyncSettings) TableName() string {
	return "tenant_sync_settings"
}
