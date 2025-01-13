package interface

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type ActionService interface {
	GetActionsForNodes(ctx context.Context, entityType model.EntityType, ids []string) (*entity.ActionEntities, error)
	CreateActionForOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string, actionFields data_fields.ActionFields) (string, error)
}
