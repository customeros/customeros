package entity

import (
	"database/sql/driver"
	"time"
)

type RawDataSource string

const (
	AirbyteSourceHubspot        RawDataSource = "hubspot"
	AirbyteSourceZendeskSupport RawDataSource = "zendesk_support"
	AirbyteSourcePipedrive      RawDataSource = "pipedrive"
	AirbyteSourceIntercom       RawDataSource = "intercom"
	OpenlineSourceSlack         RawDataSource = "slack"
)

type TenantSyncSettings struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Tenant    string    `gorm:"column:tenant;not null;uniqueIndex:uix_tenant_source_instance"`
	Source    string    `gorm:"type:string;column:source;not null;uniqueIndex:uix_tenant_source_instance"`
	Instance  string    `gorm:"column:instance;not null;default:'';uniqueIndex:uix_tenant_source_instance"`
	Enabled   bool      `gorm:"column:enabled;not null;default:false"`
}

func (rds *RawDataSource) Scan(value interface{}) error {
	*rds = RawDataSource(value.(string))
	return nil
}

func (rds RawDataSource) Value() (driver.Value, error) {
	return string(rds), nil
}

type TenantSyncSettingsList []TenantSyncSettings

func (TenantSyncSettings) TableName() string {
	return "tenant_sync_settings"
}
