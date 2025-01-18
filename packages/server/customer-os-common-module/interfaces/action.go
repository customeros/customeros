package interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type ActionService interface {
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	GetActionsForNodes(ctx context.Context, entityType model.EntityType, ids []string) (*neo4j_entity.ActionEntities, error)
	CreateActionForOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string, actionFields data_fields.ActionFields) (string, error)
}
