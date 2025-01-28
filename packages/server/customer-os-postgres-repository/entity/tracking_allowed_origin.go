package postgres_entity

import "time"

// Deprecated: TrackingAllowedOrigin is deprecated and will be removed in a future release.
type TrackingAllowedOrigin struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Tenant    string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_tenant_origin" json:"tenant" binding:"required"`
	Origin    string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_tenant_origin" json:"origin" binding:"required"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
}

func (TrackingAllowedOrigin) TableName() string {
	return "tracking_allowed_origin"
}
