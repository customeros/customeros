package entity

import (
	"fmt"
	"time"
)

type UserEntity struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
}

func (User UserEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", User.Id, User.FirstName, User.LastName)
}

type UserEntities []UserEntity
