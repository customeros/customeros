package interfaces

import (
	"context"
	common_utils "github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type AuthenticationService interface {
	SetUserService(user UserService)
	SetEmailService(email EmailService)
	IsInitialized() bool

	CreateUserInTenant(ctx context.Context, txWithPostCommit *common_utils.TxWithPostCommit, tenant string, impersonating bool, authUserId, email, firstName, lastName string) (string, error)
}
