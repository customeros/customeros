package interfaces

import (
	"context"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
)

type TenantSettingsService interface {
	GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error)
	GetTenantSettingsForTenant(ctx context.Context, tenant string) (*neo4jentity.TenantSettingsEntity, error)
	UpdateTenantSettings(ctx context.Context, dataFields data_fields.TenantSettingsFields) error

	CreateBankAccount(ctx context.Context, dataFields data_fields.BankAccountFields) (string, error)
	UpdateBankAccount(ctx context.Context, bankAccountId string, dataFields data_fields.BankAccountFields) error
	DeleteBankAccount(ctx context.Context, bankAccountId string) error

	GetTenantBillingProfiles(ctx context.Context) (*neo4jentity.TenantBillingProfileEntities, error)
	GetTenantBillingProfile(ctx context.Context, id string) (*neo4jentity.TenantBillingProfileEntity, error)
	GetDefaultTenantBillingProfile(ctx context.Context) (*neo4jentity.TenantBillingProfileEntity, error)
	CreateTenantBillingProfile(ctx context.Context, dataFields data_fields.TenantBillingProfileFields) (string, error)
	UpdateTenantBillingProfile(ctx context.Context, profileId string, dataFields data_fields.TenantBillingProfileFields) error
}
