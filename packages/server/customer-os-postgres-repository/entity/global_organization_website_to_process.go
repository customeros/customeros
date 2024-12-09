package entity

import "time"

type GlobalOrganizationWebsiteToProcess struct {
	ID          uint64    `gorm:"primary_key;autoIncrement" json:"id"`
	Website     string    `gorm:"column:website;type:varchar(255)" json:"website"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	ProcessedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"processedAt"`
	Processed   bool      `gorm:"column:processed;type:boolean" json:"processed"`
	Notes       string    `gorm:"column:notes;type:text" json:"notes"`
}

// TableName sets the name of the table for GORM
func (GlobalOrganizationWebsiteToProcess) TableName() string {
	return "global_organization_website_to_process"
}
