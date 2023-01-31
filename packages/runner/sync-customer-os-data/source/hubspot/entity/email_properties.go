package entity

import "database/sql"

type EmailProperties struct {
	AirbyteAbId         string          `gorm:"column:_airbyte_ab_id"`
	AirbyteEmailsHashid string          `gorm:"column:_airbyte_engagements_emails_hashid"`
	EmailHtml           string          `gorm:"column:hs_email_html"`
	EmailSubject        string          `gorm:"column:hs_email_subject"`
	EmailThreadId       string          `gorm:"column:hs_email_thread_id"`
	EmailDirection      string          `gorm:"column:hs_email_direction"`
	EmailFromEmail      string          `gorm:"column:hs_email_from_email"`
	EmailToEmail        string          `gorm:"column:hs_email_to_email"`
	EmailCcEmail        string          `gorm:"column:hs_email_cc_email"`
	EmailBccEmail       string          `gorm:"column:hs_email_bcc_email"`
	EmailFromFirstName  string          `gorm:"column:hs_email_from_firstname"`
	EmailFromLastName   string          `gorm:"column:hs_email_from_lastname"`
	CreatedByUserId     sql.NullFloat64 `gorm:"column:hs_created_by_user_id"`
	EmailMessageId      string          `gorm:"column:hs_email_message_id"`
}

type EmailPropertiesList []EmailProperties

func (EmailProperties) TableName() string {
	return "engagements_emails_properties"
}
