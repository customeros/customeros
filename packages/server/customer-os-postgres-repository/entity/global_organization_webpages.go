package postgres_entity

import (
	"time"

	"github.com/lib/pq"
)

type GlobalOrganizationWebpages struct {
	ID            uint64         `gorm:"primary_key;autoIncrement" json:"id"`
	PrimaryDomain string         `gorm:"column:primary_domain;type:varchar(255);NOT NULL;index:idx_global_organization_primary_domain,unique" json:"primaryDomain"`
	Url           string         `gorm:"column:url;type:text" json:"url"`
	Content       string         `gorm:"column:content;type:text" json:"content"`
	Links         pq.StringArray `gorm:"column:links;type:text[]" json:"links"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
}

func (GlobalOrganizationWebpages) TableName() string {
	return "global_organization_webpages"
}
