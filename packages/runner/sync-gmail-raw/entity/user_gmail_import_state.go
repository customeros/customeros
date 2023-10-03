package entity

import "github.com/google/uuid"

type UserGmailImportState struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantName string    `gorm:"size:255;not null"`
	Username   string    `gorm:"size:255;not null"`
	HistoryId  string    `gorm:"size:255;not null"`
}

func (UserGmailImportState) TableName() string {
	return "user_gmail_import_state"
}
