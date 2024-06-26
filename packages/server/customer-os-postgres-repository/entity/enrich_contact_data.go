package entity

import "time"

type EnrichContactData struct {
	ID        uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	ContactID string    `gorm:"column:contact_id;type:varchar(255);NOT NULL" json:"contactId"`
	App       string    `gorm:"column:app;type:varchar(255);NOT NULL" json:"app"`
	Request   string    `gorm:"column:request;type:text;DEFAULT:'';NOT NULL" json:"request"`
	Data      string    `gorm:"column:data;type:text;DEFAULT:'';NOT NULL" json:"data"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;;DEFAULT:current_timestamp" json:"updatedAt"`
}

func (EnrichContactData) TableName() string {
	return "enrich_contact_data"
}
