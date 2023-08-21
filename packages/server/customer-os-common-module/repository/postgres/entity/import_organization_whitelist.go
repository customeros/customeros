package entity

import "time"

type ImportAllowedOrganization struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	AppSource string    `gorm:"column:app_source;type:varchar(50);NOT NULL;" json:"appSource" binding:"required"`

	Tenant  string `gorm:"column:tenant;type:varchar(255);NOT NULL;index:name_domain_idx,unique" json:"tenant" binding:"required"`
	Name    string `gorm:"column:name;type:varchar(255);NOT NULL;index:name_domain_idx,unique" json:"name" binding:"required"`
	Domain  string `gorm:"column:domain;type:varchar(100);NOT NULL;index:name_domain_idx,unique" json:"domain" binding:"required"`
	Source  string `gorm:"column:source;type:varchar(100);NOT NULL;" json:"source" binding:"required"`
	Allowed bool   `gorm:"column:allowed;" json:"allowed" binding:"required"`
}

func (ImportAllowedOrganization) TableName() string {
	return "import_allowed_organization"
}
