package entity

import "time"

type Owner struct {
	Id                  string    `gorm:"column:id"`
	UserId              int64     `gorm:"column:userid"`
	AirbyteAbId         string    `gorm:"column:_airbyte_ab_id"`
	AirbyteOwnersHashid string    `gorm:"column:_airbyte_owners_hashid"`
	CreateDate          time.Time `gorm:"column:createdat"`
	UpdatedDate         time.Time `gorm:"column:updatedat"`
	FirstName           string    `gorm:"column:firstname"`
	LastName            string    `gorm:"column:lastname"`
	Email               string    `gorm:"column:email"`
}

type Owners []Owner

func (Owner) TableName() string {
	return "owners"
}
