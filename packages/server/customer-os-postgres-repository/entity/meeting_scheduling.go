package postgres_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

type MeetingScheduling struct {
	ID        string    `gorm:"type:varchar(21);primaryKey" json:"id"`
	Tenant    string    `gorm:"size:255;not null;primaryKey" json:"tenant"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updated_at"`
	Title     string    `gorm:"size:255" json:"title"`
}

func (MeetingScheduling) TableName() string {
	return "meeting_scheduling"
}

// BeforeCreate hook to ensure ID has the correct prefix and validate
func (u *MeetingScheduling) BeforeCreate(tx *gorm.DB) error {
	u.ID = utils.GenerateNanoIdWithPrefix("msch", 16)
	return nil
}

// BeforeUpdate hook to update the UpdatedAt timestamp and validate
func (u *MeetingScheduling) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}
