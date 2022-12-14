package entity

type ContactProperties struct {
	AirbyteAbId           string `gorm:"column:_airbyte_ab_id"`
	AirbyteContactsHashid string `gorm:"column:_airbyte_contacts_hashid"`
}

type ContactPropertiesList []ContactProperties

func (ContactProperties) TableName() string {
	return "contacts_properties"
}
