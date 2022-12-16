package entity

type ContactProperties struct {
	AirbyteAbId           string `gorm:"column:_airbyte_ab_id"`
	AirbyteContactsHashid string `gorm:"column:_airbyte_contacts_hashid"`
	FirstName             string `gorm:"column:firstname"`
	LastName              string `gorm:"column:lastname"`
	Email                 string `gorm:"column:email"`
	AdditionalEmails      string `gorm:"column:hs_additional_emails"`
	PhoneNumber           string `gorm:"column:phone"`
}

type ContactPropertiesList []ContactProperties

func (ContactProperties) TableName() string {
	return "contacts_properties"
}
