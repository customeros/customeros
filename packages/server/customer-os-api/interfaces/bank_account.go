package cosapi_interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type BankAccountService interface {
	GetTenantBankAccounts(ctx context.Context) (*entity.BankAccountEntities, error)
	GetTenantBankAccount(ctx context.Context, id string) (*entity.BankAccountEntity, error)
}
