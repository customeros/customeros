package entity

import (
	"github.com/jackc/pgtype"
	"time"
)

type Organization struct {
	Id                         int64        `gorm:"column:id"`
	AirbyteAbId                string       `gorm:"column:_airbyte_ab_id"`
	AirbyteOrganizationsHashid string       `gorm:"column:_airbyte_organizations_hashid"`
	CreateDate                 time.Time    `gorm:"column:created_at"`
	UpdatedDate                time.Time    `gorm:"column:updated_at"`
	DomainNames                pgtype.JSONB `gorm:"column:domain_names;type:jsonb"`
	Name                       string       `gorm:"column:name"`
	Details                    string       `gorm:"column:details"`
}

type Organizations []Organization

func (Organization) TableName() string {
	return "organizations"
}
