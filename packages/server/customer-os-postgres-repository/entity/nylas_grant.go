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

type NylasGrantResponse struct {
	RequestID string `json:"request_id"`
	Data      struct {
		ID             string   `json:"id"`
		GrantStatus    string   `json:"grant_status"`
		Provider       string   `json:"provider"`
		Scope          []string `json:"scope"`
		State          string   `json:"state"`
		Email          string   `json:"email"`
		Name           string   `json:"name"`
		IP             string   `json:"ip"`
		UserAgent      string   `json:"user_agent"`
		CreatedAt      int64    `json:"created_at"`
		UpdatedAt      int64    `json:"updated_at"`
		IDToken        string   `json:"id_token"`
		ProviderUserID string   `json:"provider_user_id"`
		Blocked        bool     `json:"blocked"`
	} `json:"data"`
}
