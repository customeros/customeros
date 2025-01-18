package interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
)

type AttachmentService interface {
	GetById(ctx context.Context, id string) (*neo4j_entity.AttachmentEntity, error)
	GetFor(ctx context.Context, entityType model.EntityType, relation *model.EntityRelation, ids []string) (*neo4j_entity.AttachmentEntities, error)
	Create(ctx context.Context, record *neo4j_entity.AttachmentEntity) (*neo4j_entity.AttachmentEntity, error)
}
