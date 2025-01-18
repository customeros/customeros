package cosapi_interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type BankAccountService interface {
	GetTenantBankAccounts(ctx context.Context) (*entity.BankAccountEntities, error)
	GetTenantBankAccount(ctx context.Context, id string) (*entity.BankAccountEntity, error)
}
