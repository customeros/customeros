package entity

import "time"

type Contact struct {
	Id                    string    `gorm:"column:id"`
	AirbyteAbId           string    `gorm:"column:_airbyte_ab_id"`
	AirbyteContactsHashid string    `gorm:"column:_airbyte_contacts_hashid"`
	CreateDate            time.Time `gorm:"column:createdat"`
	UpdatedDate           time.Time `gorm:"column:updatedat"`
}

type Contacts []Contact

func (Contact) TableName() string {
	return "contacts"
}
