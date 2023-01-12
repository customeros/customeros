package entity

type EmailProperties struct {
	AirbyteAbId         string `gorm:"column:_airbyte_ab_id"`
	AirbyteEmailsHashid string `gorm:"column:_airbyte_engagements_emails_hashid"`
	EmailHtml           string `gorm:"column:hs_email_html"`
	//OwnerId            string          `gorm:"column:hubspot_owner_id"`
	//CreatedByUserId    sql.NullFloat64 `gorm:"column:hs_created_by"`
}

type EmailPropertiesList []EmailProperties

func (EmailProperties) TableName() string {
	return "engagements_emails_properties"
}
