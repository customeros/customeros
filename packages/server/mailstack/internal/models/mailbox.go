package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Mailbox struct {
	ID        string         `gorm:"type:varchar(50);primaryKey"`
	Server    string         `gorm:"type:varchar(255);not null"`
	Port      int            `gorm:"not null"`
	Username  string         `gorm:"type:varchar(255);not null"`
	Password  string         `gorm:"type:varchar(255);not null"`
	Folders   pq.StringArray `gorm:"type:text;not null"`
	TLS       bool           `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
