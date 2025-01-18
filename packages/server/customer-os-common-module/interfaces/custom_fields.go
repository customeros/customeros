package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
)

type CustomFieldTemplateService interface {
	GetAll(ctx context.Context) (*entity.CustomFieldTemplateEntities, error)
	GetById(ctx context.Context, customFieldTemplateId string) (*entity.CustomFieldTemplateEntity, error)
	Save(ctx context.Context, id *string, input repository.CustomFieldTemplateSaveFields) (string, error)
	Delete(ctx context.Context, customFieldTemplateId string) error
}
