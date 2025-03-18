package interfaces

import "context"

type RegistrationService interface {
	SetContactService(contact ContactService)
	SetEmailService(email EmailService)
	SetFlowService(flow FlowService)
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error
	SetupTestUser(ctx context.Context) (*TestUserSetup, error)
	SetupTestMailbox(ctx context.Context, tenant string, testUser *TestUserSetup) error
	ConfigureTestMailbox(ctx context.Context) (*TestUserSetup, error)
}

type TestUserSetup struct {
	UserId         string
	MailboxAddress string
}
