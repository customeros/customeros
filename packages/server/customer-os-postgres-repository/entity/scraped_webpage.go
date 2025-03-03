package postgres_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/lib/pq"
)

type ScrapedWebpage struct {
	ID            uint64                    `gorm:"primary_key;autoIncrement" json:"id"`
	PrimaryDomain string                    `gorm:"column:primary_domain;type:varchar(255);NOT NULL;index:idx_global_organization_primary_domain,unique" json:"primaryDomain"`
	Url           string                    `gorm:"column:url;type:text;uniqueIndex:idx_unique_url" json:"url"`
	Content       string                    `gorm:"column:content;type:text" json:"content"`
	Links         pq.StringArray            `gorm:"column:links;type:text[]" json:"links"`
	ContentStage  enum.CustomerJourneyStage `gorm:"column:content_stage;type:varchar(55)" json:"contentStage"`
	Category      enum.WebpageCategory      `gorm:"column:category;type:varchar(55)" json:"category"`
	Topics        pq.StringArray            `gorm:"column:topics;type:text[]" json:"topics"`
	CreatedAt     time.Time                 `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time                 `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Error         string                    `gorm:"column:error;type:varchar(255)" json:"error"`
}

func (ScrapedWebpage) TableName() string {
	return "scraped_webpages"
}
