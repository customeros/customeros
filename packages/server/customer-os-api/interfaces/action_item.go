package cosapi_interfaces

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type ActionItemService interface {
	GetActionItemsForNodes(ctx context.Context, linkedWith repository.LinkedWith, ids []string) (*entity.ActionItemEntities, error)

	MapDbNodeToActionItemEntity(node dbtype.Node) *entity.ActionItemEntity
}
