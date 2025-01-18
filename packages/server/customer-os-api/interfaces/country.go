package cosapi_interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type CountryService interface {
	GetCountriesForPhoneNumbers(ctx context.Context, ids []string) (*entity.CountryEntities, error)
}
