package entity

import "database/sql"

type ContactProperties struct {
	AirbyteAbId              string          `gorm:"column:_airbyte_ab_id"`
	AirbyteContactsHashid    string          `gorm:"column:_airbyte_contacts_hashid"`
	FirstName                string          `gorm:"column:firstname"`
	LastName                 string          `gorm:"column:lastname"`
	Email                    string          `gorm:"column:email"`
	AdditionalEmails         string          `gorm:"column:hs_additional_emails"`
	PhoneNumber              string          `gorm:"column:phone"`
	PrimaryCompanyExternalId sql.NullFloat64 `gorm:"column:associatedcompanyid"`
	JobTitle                 string          `gorm:"column:jobtitle"`
	OwnerId                  string          `gorm:"column:hubspot_owner_id"`
	LifecycleStage           string          `gorm:"column:lifecyclestage"`
}

type ContactPropertiesList []ContactProperties

func (ContactProperties) TableName() string {
	return "contacts_properties"
}
