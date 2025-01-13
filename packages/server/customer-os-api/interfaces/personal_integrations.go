package cosapi_interfaces

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

type PersonalIntegrationsService interface {
	GetPersonalIntegration(tenantName, email, integration string) (*entity.PersonalIntegration, error)
	SavePersonalIntegration(entity.PersonalIntegration) (*entity.PersonalIntegration, error)
	GetPersonalIntegrations(tenantName, email string) ([]*entity.PersonalIntegration, error)
}
