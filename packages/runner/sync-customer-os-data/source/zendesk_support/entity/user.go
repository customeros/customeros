package entity

import (
	"time"
)

type User struct {
	Id                 int64     `gorm:"column:id"`
	AirbyteAbId        string    `gorm:"column:_airbyte_ab_id"`
	AirbyteUsersHashid string    `gorm:"column:_airbyte_users_hashid"`
	CreateDate         time.Time `gorm:"column:created_at"`
	UpdatedDate        time.Time `gorm:"column:updated_at"`
	Name               string    `gorm:"column:name"`
	Email              string    `gorm:"column:email"`
	Phone              string    `gorm:"column:phone"`
	Role               string    `gorm:"column:role"`
}

type Users []User

func (User) TableName() string {
	return "users"
}

func (user User) IsEndUser() bool {
	return user.Role == "end-user"
}
