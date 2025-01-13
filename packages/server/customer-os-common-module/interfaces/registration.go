package interfaces

import "context"

type RegistrationService interface {
	PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error
	ConfigureTestMailbox(ctx context.Context) (*TestUserSetup, error)
	ConfigureDefaultFlowData(ctx context.Context, testUser *TestUserSetup) error
	CreatePostmarkServer(ctx context.Context) error
}

type TestUserSetup struct {
	UserId         string
	MailboxAddress string
}
