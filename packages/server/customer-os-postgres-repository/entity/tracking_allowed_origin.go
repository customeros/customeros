package entity

import "time"

type TrackingAllowedOrigin struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL;index:name_domain_idx,unique" json:"tenant" binding:"required"`
	Origin    string    `gorm:"column:origin;type:varchar(255);NOT NULL;index:name_domain_idx,unique" json:"name" binding:"required"`
}

func (TrackingAllowedOrigin) TableName() string {
	return "tracking_allowed_origin"
}
