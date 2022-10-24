package entity

import (
	"time"
)

type ApplicationUniqueIdentifier struct {
	AppId       string
	TrackerName string
	Tenant      string
}

type ApplicationEntity struct {
	ID          string    `gorm:"primary_key" json:"id"`
	Platform    string    `gorm:"column:platform;type:varchar(255);NOT NULL" json:"platform" binding:"required"`
	AppId       string    `gorm:"column:app_id;type:varchar(255);NOT NULL" json:"name" binding:"required"`
	TrackerName string    `gorm:"column:name_tracker;type:varchar(128);NOT NULL" json:"trackerName" binding:"required"`
	UpdatedOn   time.Time `gorm:"column:updated_on;NOT NULL" json:"updatedOn" binding:"required"`
	Tenant      string    `gorm:"column:tenant;type:varchar(64);NOT NULL" json:"tenant" binding:"required"`
}

type ApplicationEntities []ApplicationEntity

func (ApplicationEntity) TableName() string {
	return "derived.app_info"
}
