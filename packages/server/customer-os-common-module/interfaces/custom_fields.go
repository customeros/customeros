package interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
)

type CustomFieldTemplateService interface {
	GetAll(ctx context.Context) (*neo4j_entity.CustomFieldTemplateEntities, error)
	GetById(ctx context.Context, customFieldTemplateId string) (*neo4j_entity.CustomFieldTemplateEntity, error)
	Save(ctx context.Context, id *string, input neo4j_repository.CustomFieldTemplateSaveFields) (string, error)
	Delete(ctx context.Context, customFieldTemplateId string) error
}
