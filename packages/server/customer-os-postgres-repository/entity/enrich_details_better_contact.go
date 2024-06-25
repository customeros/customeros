package entity

import "time"

type EnrichDetailsBetterContact struct {
	ID        uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	ContactID string    `gorm:"column:contact_id;type:varchar(255);NOT NULL" json:"contactId"`
	RequestID string    `gorm:"column:request_id;type:varchar(255);NOT NULL" json:"requestId"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;" json:"updatedAt"`
	Request   string    `gorm:"column:request;type:text;DEFAULT:'';NOT NULL" json:"request"`
	Response  string    `gorm:"column:response;type:text;DEFAULT:'';NOT NULL" json:"response"`
}

func (EnrichDetailsBetterContact) TableName() string {
	return "enrich_details_better_contact"
}
