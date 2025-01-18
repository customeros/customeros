package cosapi_interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type CountryService interface {
	GetCountriesForPhoneNumbers(ctx context.Context, ids []string) (*entity.CountryEntities, error)
}
