package interfaces

import (
	"context"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type TaskService interface {
	IsInitialized() bool

	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, taskFields data_fields.TaskFields) (string, error)
	HideAll(ctx context.Context, ids []string) error
	Hide(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id string) error
	GetById(ctx context.Context, id string) (*neo4jentity.TaskEntity, error)
	GetAllByIds(ctx context.Context, ids []string) (*neo4jentity.TaskEntities, error)
	GetTasksForOpportunities(ctx context.Context, opportunityIds []string) (*neo4jentity.TaskEntities, error)
}
