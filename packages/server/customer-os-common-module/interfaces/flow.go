package interfaces

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
)

type FlowService interface {
	SetFlowExecutionService(fe FlowExecutionService)
	IsInitialized() bool

	FlowGetList(ctx context.Context) (*neo4jentity.FlowEntities, error)
	FlowGetById(ctx context.Context, id string) (*neo4jentity.FlowEntity, error)
	FlowGetByActionId(ctx context.Context, flowActionId string) (*neo4jentity.FlowEntity, error)
	FlowGetByParticipantId(ctx context.Context, tx *neo4j.ManagedTransaction, flowParticipantId string) (*neo4jentity.FlowEntity, error)
	FlowsGetListWithParticipant(ctx context.Context, entityIds []string, entityType model.EntityType) (*neo4jentity.FlowEntities, error)
	FlowsGetListWithSender(ctx context.Context, senderIds []string) (*neo4jentity.FlowEntities, error)
	FlowMerge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *neo4jentity.FlowEntity) (*neo4jentity.FlowEntity, error)
	FlowOn(ctx context.Context, id string) (*neo4jentity.FlowEntity, error)
	FlowOff(ctx context.Context, id string) (*neo4jentity.FlowEntity, error)
	FlowArchive(ctx context.Context, id string) (*neo4jentity.FlowEntity, error)

	FlowActionGetStart(ctx context.Context, flowId string) (*neo4jentity.FlowActionEntity, error)
	FlowActionGetNext(ctx context.Context, actionId string) ([]*neo4jentity.FlowActionEntity, error)
	FlowActionGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowActionEntities, error)
	FlowActionGetById(ctx context.Context, id string) (*neo4jentity.FlowActionEntity, error)

	FlowParticipantGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowParticipantEntities, error)
	FlowParticipantById(ctx context.Context, flowParticipantId string) (*neo4jentity.FlowParticipantEntity, error)
	FlowParticipantByEntity(ctx context.Context, flowId, entityId string, entityType model.EntityType) (*neo4jentity.FlowParticipantEntity, error)
	FlowParticipantAdd(ctx context.Context, flowId, entityId string, entityType model.EntityType) (*neo4jentity.FlowParticipantEntity, error)
	FlowParticipantDelete(ctx context.Context, flowParticipantId string) error

	FlowSenderGetList(ctx context.Context, flowIds []string) (*neo4jentity.FlowSenderEntities, error)
	FlowSenderGetById(ctx context.Context, id string) (*neo4jentity.FlowSenderEntity, error)
	FlowSenderMerge(ctx context.Context, flowId string, input *neo4jentity.FlowSenderEntity) (*neo4jentity.FlowSenderEntity, error)
	FlowSenderDelete(ctx context.Context, flowSenderId string) error
}
