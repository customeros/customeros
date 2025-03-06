package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type GlobalContactService interface {
	SaveContact(ctx context.Context, contact *postgres_entity.GlobalContact) error
}
