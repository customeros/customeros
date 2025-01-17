package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
)

type PhoneNumberService interface {
	Merge(ctx context.Context, phoneNumber string, source entity.DataSource) (string, error)
	UpdatePhoneNumberFor(ctx context.Context, entityType model.EntityType, entityId string, phoneId string, label *string, primary *bool) error
	DetachFromEntityByPhoneNumber(ctx context.Context, entityType model.EntityType, entityId, phoneNumber string) (bool, error)
	DetachFromEntityById(ctx context.Context, entityType model.EntityType, entityId, phoneNumberId string) (bool, error)
	GetAllForEntityTypeByIds(ctx context.Context, entityType model.EntityType, ids []string) (*entity.PhoneNumberEntities, error)
	GetById(ctx context.Context, phoneNumberId string) (*entity.PhoneNumberEntity, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.PhoneNumberEntity, error)
}
