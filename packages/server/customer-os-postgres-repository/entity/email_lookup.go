package entity

import "time"

type EmailLookupType string

const (
	EmailLookupTypeLink     EmailLookupType = "Link"
	EmailLookupTypeSpyPixel EmailLookupType = "SpyPixel"
)

type EmailLookup struct {
	ID          string          `gorm:"column:id,primary_key;type:varchar(64);" json:"id"`
	Tenant      string          `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt   time.Time       `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time       `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	MessageId   string          `gorm:"column:message_id;type:varchar(64);NOT NULL" json:"messageId"`
	LinkId      string          `gorm:"column:link_id;type:varchar(64);NOT NULL" json:"linkId"`
	RedirectUrl string          `gorm:"column:redirect_url;type:varchar(255);NOT NULL" json:"redirectUrl"`
	Campaign    string          `gorm:"column:campaign;type:varchar(255);NOT NULL" json:"campaign"`
	Type        EmailLookupType `gorm:"column:type;type:varchar(255);NOT NULL" json:"type"`
}

func (EmailLookup) TableName() string {
	return "email_lookup"
}
