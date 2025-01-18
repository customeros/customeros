package interfaces

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type FlowExecutionService interface {
	SetEmailService(email EmailService)
	SetFlowService(flow FlowService)
	SetOrganizationService(org OrganizationService)
	SetSocialService(social SocialService)
	IsInitialized() bool

	GetFlowActionExecutionById(ctx context.Context, flowActionExecution string) (*entity.FlowActionExecutionEntity, error)
	GetFlowExecutionSettingsForEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType) (*entity.FlowExecutionSettingsEntity, error)
	GetFlowRequirements(ctx context.Context, flowId string) (*FlowComputeParticipantsRequirementsInput, error)
	GetFlowActionExecutionsForParticipants(ctx context.Context, flowParticipantIds []string) (*entity.FlowActionExecutionEntities, error)
	GetFlowActionExecutionsForParticipant(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType) ([]*entity.FlowActionExecutionEntity, error)
	GetFlowActionExecutionsForParticipantWithActionType(ctx context.Context, entityId string, entityType model.EntityType, actionType entity.FlowActionType) ([]*entity.FlowActionExecutionEntity, error)
	UpdateAllParticipantsFlowRequirements(ctx context.Context, flowId string) error
	UpdateParticipantFlowRequirements(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, participant *entity.FlowParticipantEntity, requirements *FlowComputeParticipantsRequirementsInput) error
	ScheduleFlow(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, flowId string, flowParticipant *entity.FlowParticipantEntity) error
	ProcessActionExecution(ctx context.Context, scheduledActionExecution *entity.FlowActionExecutionEntity) error
	ReplacePlaceholder(input, variableName, value string) string
}

type FlowComputeParticipantsRequirementsInput struct {
	PrimaryEmailRequired      bool `json:"primaryEmailRequired"`
	LinkedInSocialUrlRequired bool `json:"linkedInSocialUrlRequired"`
}
