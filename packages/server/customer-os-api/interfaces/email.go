package cosapi_interfaces

import (
	"context"

	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type EmailService interface {
	GetAllFor(ctx context.Context, entityType commonModel.EntityType, entityId string) (*neo4jentity.EmailEntities, error)
	GetAllForEntityTypeByIds(ctx context.Context, entityType commonModel.EntityType, entityIds []string) (*neo4jentity.EmailEntities, error)
	GetById(ctx context.Context, emailId string) (*neo4jentity.EmailEntity, error)
	GetByEmailAddress(ctx context.Context, email string) (*neo4jentity.EmailEntity, error)
}
