package cosapi_interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type CountryService interface {
	GetCountriesForPhoneNumbers(ctx context.Context, ids []string) (*neo4j_entity.CountryEntities, error)
}
