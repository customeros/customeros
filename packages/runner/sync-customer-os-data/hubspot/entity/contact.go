package entity

import "time"

type Contact struct {
	Id                    string    `gorm:"column:id"`
	AirbyteAbId           string    `gorm:"column:_airbyte_ab_id"`
	AirbyteContactsHashid string    `gorm:"column:_airbyte_contacts_hashid"`
	UpdatedDate           time.Time `gorm:"column:updatedat"`
	SyncedToCustomerOs    bool      `gorm:"column:synced_to_customer_os"`
}

type Contacts []Contact

func (Contact) TableName() string {
	return "contacts"
}
