package postgres_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

type NylasAccount struct {
	ID             string    `gorm:"primaryKey;type:varchar(21)" json:"id"`
	Tenant         string    `gorm:"column:tenant;type:varchar(255);NOT NULL;index:idx_nylas_account_tenant_email,unique" json:"tenant"`
	Email          string    `gorm:"column:email;type:varchar(255);NOT NULL;index:idx_nylas_account_tenant_email,unique" json:"email"`
	NylasAccountId string    `gorm:"column:nylas_account_id;type:varchar(255);NOT NULL" json:"nylasAccountId"`
	Provider       string    `gorm:"column:provider;type:varchar(50);NOT NULL;index:idx_nylas_account_tenant_email,unique" json:"provider"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
}

func (NylasAccount) TableName() string {
	return "nylas_account"
}

func (n *NylasAccount) BeforeCreate(tx *gorm.DB) error {
	n.ID = utils.GenerateNanoIdWithPrefix("nylas", 16)
	return nil
}
