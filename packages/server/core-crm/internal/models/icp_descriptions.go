package models

import "time"

type ICPDescription struct {
	ID                      string     `gorm:"column:id;primaryKey;type:varchar(30);"`
	Tenant                  string     `gorm:"column:tenant;type:varchar(255);index;not null"`
	Profile                 string     `gorm:"column:profile;type:text;not null"`
	QualifyingAttributes    string     `gorm:"column:qualifying_attributes;type:text;not null"`
	DisqualifyingAttributes string     `gorm:"column:disqualifying_attributes;type:text;not null"`
	CreatedAt               time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt               *time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ICPDescription) TableName() string {
	return "icp_descriptions"
}
