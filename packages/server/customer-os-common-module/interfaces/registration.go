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
}

type TestUserSetup struct {
	UserId         string
	MailboxAddress string
}
