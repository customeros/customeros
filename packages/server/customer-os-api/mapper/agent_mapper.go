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
		Icon:         entity.Icon,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    utils.IfNotNilTimeWithDefault(entity.UpdatedAt, entity.CreatedAt),
		Type:         enummapper.MapAgentTypeToModel(entity.Type),
		Color:        entity.Color,
		Goal:         entity.Goal,
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
	return &agentModel
}
