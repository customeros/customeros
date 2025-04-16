package postgres_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

type NylasGrant struct {
	ID                   string    `gorm:"primaryKey;type:varchar(21)" json:"id"`
	Tenant               string    `gorm:"column:tenant;type:varchar(255);NOT NULL;index:idx_nylas_account_tenant_email,unique" json:"tenant"`
	Email                string    `gorm:"column:email;type:varchar(255);NOT NULL;index:idx_nylas_account_tenant_email,unique" json:"email"`
	NylasGrantId         string    `gorm:"column:nylas_grant_id;type:varchar(255);NOT NULL" json:"nylasGrantId"`
	NylasProvider        string    `gorm:"column:nylas_provider;type:varchar(50);NOT NULL;index:idx_nylas_account_tenant_email,unique" json:"provider"`
	CreatedAt            time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt            time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	NylasConnectResponse string    `gorm:"column:nylas_connect_response;type:text" json:"nylasConnectResponse"`
}

func (NylasGrant) TableName() string {
	return "nylas_grant"
}

func (n *NylasGrant) BeforeCreate(tx *gorm.DB) error {
	n.ID = utils.GenerateNanoIdWithPrefix("nylas", 16)
	return nil
}
