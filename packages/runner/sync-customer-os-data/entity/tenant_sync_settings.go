package entity

import (
	"database/sql/driver"
	"time"
)

type AirbyteSource string

const (
	AirbyteSourceHubspot        AirbyteSource = "hubspot"
	AirbyteSourceZendeskSupport AirbyteSource = "zendesk_support"
	AirbyteSourcePipedrive      AirbyteSource = "pipedrive"
)

type TenantSyncSettings struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Tenant    string    `gorm:"column:tenant;not null;uniqueIndex:uix_tenant_source_instance"`
	Source    string    `gorm:"type:string;column:source;not null;uniqueIndex:uix_tenant_source_instance"`
	Instance  string    `gorm:"column:instance;not null;default:'';uniqueIndex:uix_tenant_source_instance"`
	Enabled   bool      `gorm:"column:enabled;not null;default:false"`
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
