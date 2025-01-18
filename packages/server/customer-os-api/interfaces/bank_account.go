package cosapi_interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type BankAccountService interface {
	GetTenantBankAccounts(ctx context.Context) (*neo4j_entity.BankAccountEntities, error)
	GetTenantBankAccount(ctx context.Context, id string) (*neo4j_entity.BankAccountEntity, error)
}
