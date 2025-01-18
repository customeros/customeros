package interfaces

import "context"

type RegistrationService interface {
	SetContactService(contact ContactService)
	SetEmailService(email EmailService)
	SetFlowService(flow FlowService)
	SetMailboxService(mailbox MailboxService)
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error
	ConfigureTestMailbox(ctx context.Context) (*TestUserSetup, error)
	ConfigureDefaultFlowData(ctx context.Context, testUser *TestUserSetup) error
	CreatePostmarkServer(ctx context.Context) error
}

type TestUserSetup struct {
	UserId         string
	MailboxAddress string
}
