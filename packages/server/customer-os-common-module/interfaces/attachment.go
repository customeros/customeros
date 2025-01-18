package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
)

type AttachmentService interface {
	GetById(ctx context.Context, id string) (*entity.AttachmentEntity, error)
	GetFor(ctx context.Context, entityType model.EntityType, relation *model.EntityRelation, ids []string) (*entity.AttachmentEntities, error)
	Create(ctx context.Context, record *entity.AttachmentEntity) (*entity.AttachmentEntity, error)
}
