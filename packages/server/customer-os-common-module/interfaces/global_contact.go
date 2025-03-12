package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type GlobalContactService interface {
	SaveContact(ctx context.Context, contact *postgres_entity.GlobalContact) error
	SetWorkEmail(ctx context.Context, id uint64, workEmail string) error
	GetGlobalContactsByLinkedIn(ctx context.Context, linkedIn string) ([]*postgres_entity.GlobalContact, error)
}
