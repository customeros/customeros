package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
)

type ExternalSystemService interface {
	MergeExternalSystem(ctx context.Context, tenant, externalSystem string) error
	SetPrimaryExternalId(ctx context.Context, externalSystem, externalId string, linkWith common_srv.LinkWith) error
	GetPrimaryExternalId(ctx context.Context, externalSystem, linkedWithId string, linkedWithEntityType model.EntityType) (string, error)
	GetExternalSystemsForEntities(ctx context.Context, ids []string, entityType model.EntityType) (*entity.ExternalSystemEntities, error)
}
