package cosapi_interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

// TODO deprecate and remove
type CustomFieldService interface {
	MergeAndUpdateCustomFieldsForContact(ctx context.Context, contactId string, customFields *entity.CustomFieldEntities) error
	MergeCustomFieldToContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)
	UpdateCustomFieldForContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)
	DeleteByNameFromContact(ctx context.Context, contactId, fieldName string) (bool, error)
	DeleteByIdFromContact(ctx context.Context, contactId, fieldId string) (bool, error)
	GetCustomFields(ctx context.Context, obj *model.CustomFieldEntityType) (*entity.CustomFieldEntities, error)
}
