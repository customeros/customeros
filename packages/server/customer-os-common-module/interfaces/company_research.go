package interfaces

import "context"

type CompanyResearch interface {
	GetIdealCustomerProfile(ctx context.Context) (string, error)
}
