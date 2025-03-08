package models

import (
	"time"

	"gorm.io/gorm"
)

type Mailbox struct {
	ID        string   `gorm:"type:varchar(50);primaryKey"`
	Server    string   `gorm:"type:varchar(255);not null"`
	Port      int      `gorm:"not null"`
	Username  string   `gorm:"type:varchar(255);not null"`
	Password  string   `gorm:"type:varchar(255);not null"`
	Folders   []string `gorm:"type:jsonb;not null"`
	TLS       bool     `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
