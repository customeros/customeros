package entity

import "time"

type EmailValidationRecord struct {
	ID        uint64    `gorm:"primary_key;autoIncrement" json:"id"`
	RequestID string    `gorm:"column:request_id;type:varchar(255);NOT NULL;index:idx_email_request_id,unique" json:"requestId"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	Email     string    `gorm:"column:email;type:varchar(255);NOT NULL;index:idx_email_request_id,unique" json:"email"`
	Priority  int       `gorm:"column:priority;type:int;DEFAULT:0" json:"priority"`
	Data      string    `gorm:"column:data;type:text" json:"data"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;" json:"updatedAt"`
}

func (EmailValidationRecord) TableName() string {
	return "email_validation_record"
}
