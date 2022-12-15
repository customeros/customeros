package entity

import "time"

type Contact struct {
	Id                    string    `gorm:"column:id"`
	AirbyteAbId           string    `gorm:"column:_airbyte_ab_id"`
	AirbyteContactsHashid string    `gorm:"column:_airbyte_contacts_hashid"`
	CreateDate            time.Time `gorm:"column:createdat"`
	UpdatedDate           time.Time `gorm:"column:updatedat"`
	SyncedToCustomerOs    bool      `gorm:"column:synced_to_customer_os"`
	SyncAttempt           int       `gorm:"column:synced_to_customer_os_attempt"`
	SyncedAt              time.Time `gorm:"column:synced_to_customer_os_at"`
}

type Contacts []Contact

func (Contact) TableName() string {
	return "contacts"
}
