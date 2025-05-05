package models

import "github.com/customeros/customeros/packages/server/core-crm/internal/enum"

type Company struct {
	ID               string                    `gorm:"column:id;primaryKey;type:varchar(30);"`
	Tenant           string                    `gorm:"column:tenant;type:varchar(255);index;not null"`
	CustomID         string                    `gorm:"column:custom_id;type:varchar(255);index;not null"`
	GlobalOrgID      string                    `gorm:"column:global_org_id;type:varchar(255);index;not null"`
	ICPFit           string                    `gorm:"column:country;type:varchar(255)"`
	ICPFitUpdatedAt  string                    `gorm:"column:tenant;type:varchar(255);index;not null"`
	Stage            enum.CustomerJourneyStage `gorm:"column:stage;type:varchar(50);index;default:'target'"`
	IsCustomer       bool                      `gorm:"column:stage;type:varchar(50);index;default:'target'"`
	IsFormerCustomer bool                      `gorm:"column:stage;type:varchar(50);index;default:'target'"`
}
