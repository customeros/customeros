package entity

import "time"

type EmailLookupType string

const (
	EmailLookupTypeLink        EmailLookupType = "Link"
	EmailLookupTypeSpyPixel    EmailLookupType = "SpyPixel"
	EmailLookupTypeUnsubscribe EmailLookupType = "Unsubscribe"
)

type EmailLookup struct {
	ID             string          `gorm:"column:id;primary_key;type:varchar(64);" json:"id"`
	Tenant         string          `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt      time.Time       `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt      time.Time       `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	TrackerDomain  string          `gorm:"column:tracker_domain;type:varchar(255);" json:"trackerDomain"`
	MessageId      string          `gorm:"column:message_id;type:varchar(64);NOT NULL" json:"messageId"`
	LinkId         string          `gorm:"column:link_id;type:varchar(64);NOT NULL" json:"linkId"`
	RedirectUrl    string          `gorm:"column:redirect_url;type:varchar(255);NOT NULL" json:"redirectUrl"`
	Campaign       string          `gorm:"column:campaign;type:varchar(255);NOT NULL" json:"campaign"`
	Type           EmailLookupType `gorm:"column:type;type:varchar(32);NOT NULL" json:"type"`
	RecipientId    string          `gorm:"column:recipient_id;type:varchar(255);" json:"recipientId"`
	TrackOpens     bool            `gorm:"column:track_opens;type:boolean;NOT NULL" json:"trackOpens"`
	TrackClicks    bool            `gorm:"column:track_clicks;type:boolean;NOT NULL" json:"trackClicks"`
	UnsubscribeUrl string          `gorm:"column:unsubscribe_url;type:varchar(255);" json:"unsubscribeUrl"`
}

func (EmailLookup) TableName() string {
	return "email_lookup"
}
