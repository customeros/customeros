package interfaces

import "context"

type RegistrationService interface {
	PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error
	ConfigureTestMailbox(ctx context.Context) (*testUserSetup, error)
	ConfigureDefaultFlowData(ctx context.Context, testUser *testUserSetup) error
	CreatePostmarkServer(ctx context.Context) error
}

type testUserSetup struct {
	userId         string
	mailboxAddress string
}
