package entity

import "time"

type CacheEmailTrueinbox struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	Email     string    `gorm:"column:email;type:varchar(255);NOT NULL" json:"email"`
	Result    string    `gorm:"column:result;type:varchar(255)" json:"result"`
	Data      string    `gorm:"column:data;type:text" json:"data"`
}

func (CacheEmailTrueinbox) TableName() string {
	return "cache_email_trueinbox"
}

type TrueInboxResponseBody struct {
	Email            string `json:"email"`
	Status           string `json:"status"`
	Result           string `json:"result"`
	ConfidenceScore  int    `json:"confidenceScore"`
	SmtpProvider     string `json:"smtpProvider"`
	MailDisposable   bool   `json:"mailDisposable"`
	MailAcceptAll    bool   `json:"mailAcceptAll"`
	Free             bool   `json:"free"`
	TotalCredits     int    `json:"total_credits"`
	CreditsUsed      int    `json:"credits_used"`
	CreditsRemaining int    `json:"credits_remaining"`
}
