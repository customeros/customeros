package entity

import (
	"time"
)

type UserAsOrganization struct {
	Id                   int64     `gorm:"column:id"`
	AirbyteAbId          string    `gorm:"column:_airbyte_ab_id"`
	AirbyteUsersHashid   string    `gorm:"column:_airbyte_users_hashid"`
	CreateDate           time.Time `gorm:"column:created_at"`
	UpdatedDate          time.Time `gorm:"column:updated_at"`
	Name                 string    `gorm:"column:name"`
	Email                string    `gorm:"column:email"`
	Phone                string    `gorm:"column:phone"`
	Url                  string    `gorm:"column:url"`
	Details              string    `gorm:"column:details"`
	Notes                string    `gorm:"column:notes"`
	ParentOrganizationId int64     `gorm:"column:organization_id"`
}

type UsersAsOrganizations []UserAsOrganization

func (UserAsOrganization) TableName() string {
	return "users"
}
