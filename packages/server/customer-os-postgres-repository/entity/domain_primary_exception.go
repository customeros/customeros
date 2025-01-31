package postgres_entity

import (
	"time"
)

type DomainPrimaryException struct {
	ID        uint64    `gorm:"primary_key;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	Domain    string    `gorm:"column:domain;type:varchar(255);not null" json:"domain"`
}

func (DomainPrimaryException) TableName() string {
	return "domain_primary_exception"
}
