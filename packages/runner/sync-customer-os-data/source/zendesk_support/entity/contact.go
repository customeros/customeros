package entity

import (
	"time"
)

type Contact struct {
	Id                 int64     `gorm:"column:id"`
	AirbyteAbId        string    `gorm:"column:_airbyte_ab_id"`
	AirbyteUsersHashid string    `gorm:"column:_airbyte_users_hashid"`
	CreateDate         time.Time `gorm:"column:created_at"`
	UpdatedDate        time.Time `gorm:"column:updated_at"`
	Name               string    `gorm:"column:name"`
	Email              string    `gorm:"column:email"`
	Phone              string    `gorm:"column:phone"`
	Url                string    `gorm:"column:url"`
	Notes              string    `gorm:"column:notes"`
	Details            string    `gorm:"column:details"`
	OrganizationId     int64     `gorm:"column:organization_id"`
}

type Contacts []Contact

func (Contact) TableName() string {
	return "users"
}
