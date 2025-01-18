package cosapi_interfaces

import postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

type PersonalIntegrationsService interface {
	GetPersonalIntegration(tenantName, email, integration string) (*postgres_entity.PersonalIntegration, error)
	SavePersonalIntegration(postgres_entity.PersonalIntegration) (*postgres_entity.PersonalIntegration, error)
	GetPersonalIntegrations(tenantName, email string) ([]*postgres_entity.PersonalIntegration, error)
}
