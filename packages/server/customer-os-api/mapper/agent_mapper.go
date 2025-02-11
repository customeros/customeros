package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	enummapper "github.com/customeros/customeros/packages/server/customer-os-api/mapper/enum"
)

func MapAgentToModel(entity *postgresEntity.Agent) *model.Agent {
	if entity == nil {
		return nil
	}
	agentModel := model.Agent{
		ID:           entity.ID,
		Name:         entity.Name,
		Scope:        model.AgentScope(entity.Scope),
		Icon:         entity.Icon,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    utils.IfNotNilTimeWithDefault(entity.UpdatedAt, entity.CreatedAt),
		Type:         enummapper.MapAgentTypeToModel(entity.Type),
		Color:        entity.Color,
		GoalType:     entity.Goal.String(),
		Goal:         enummapper.MapAgentGoalName(entity.Goal),
		IsActive:     entity.IsActive,
		IsConfigured: entity.Configured,
		Visible:      entity.VisibleInUI,
		FlowID:       utils.StringPtr(entity.FlowID),
	}
	for _, capability := range entity.Capabilities {
		agentModel.Capabilities = append(agentModel.Capabilities, &model.Capability{
			ID:     capability.ID,
			Name:   capability.Name,
			Type:   enummapper.MapAgentCapabilityTypeToModel(capability.Type),
			Errors: utils.StringPtrNillable(capability.Error),
			Active: capability.Active,
			Config: capability.GetConfigString(),
		})
	}
	for _, listener := range entity.Listeners {
		agentModel.Listeners = append(agentModel.Listeners, &model.AgentListener{
			ID:     listener.ID,
			Name:   listener.Name,
			Type:   enummapper.MapAgentListenerTypeToModel(listener.Type),
			Errors: utils.StringPtrNillable(listener.Error),
			Active: listener.Active,
			Config: listener.GetConfigString(),
		})
	}
	return &agentModel
}
