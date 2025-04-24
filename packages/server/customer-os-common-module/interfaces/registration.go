package interfaces

import "context"

type RegistrationService interface {
	IsInitialized() bool

	InitialTenantSetup(ctx context.Context, loggedInUserEmail string) error
	ProvideAccessToPlatformOwners(ctx context.Context, tenant string) error
	SetupTestUser(ctx context.Context) (*TestUserSetup, error)
	SetupTestMailbox(ctx context.Context, tenant string, testUser *TestUserSetup) error
	ConfigureTestMailbox(ctx context.Context) (*TestUserSetup, error)
}

type TestUserSetup struct {
	UserId         string
	MailboxAddress string
}
