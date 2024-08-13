package entity

import (
	"time"
)

type OrganizationWebsiteHostingPlatform struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UrlPattern string    `gorm:"size:255;NOT NULL"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (OrganizationWebsiteHostingPlatform) TableName() string {
	return "organization_website_hosting_platform"
}
