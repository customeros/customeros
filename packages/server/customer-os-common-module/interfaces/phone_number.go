package interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
)

type PhoneNumberService interface {
	Merge(ctx context.Context, phoneNumber string, source neo4j_entity.DataSource) (string, error)
	UpdatePhoneNumberFor(ctx context.Context, entityType model.EntityType, entityId string, phoneId string, label *string, primary *bool) error
	DetachFromEntityByPhoneNumber(ctx context.Context, entityType model.EntityType, entityId, phoneNumber string) (bool, error)
	DetachFromEntityById(ctx context.Context, entityType model.EntityType, entityId, phoneNumberId string) (bool, error)
	GetAllForEntityTypeByIds(ctx context.Context, entityType model.EntityType, ids []string) (*neo4j_entity.PhoneNumberEntities, error)
	GetById(ctx context.Context, phoneNumberId string) (*neo4j_entity.PhoneNumberEntity, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*neo4j_entity.PhoneNumberEntity, error)
}
