package interfaces

import "context"

type CompanyResearch interface {
	GenerateIdealCustomerProfile(ctx context.Context, tenantDomain string, trainingWebsites []string) (string, error)
	GenerateCompanyBrief(ctx context.Context, domains []string) (string, error)
}
