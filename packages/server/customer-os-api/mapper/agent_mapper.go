package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	enummapper "github.com/customeros/customeros/packages/server/customer-os-api/mapper/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

func MapAgentToModel(entity *postgresEntity.Agents) *model.Agent {
	if entity == nil {
		return nil
	}
	agentModel := model.Agent{
		ID:        entity.ID,
		Name:      entity.Name,
		Icon:      entity.Icon,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: utils.IfNotNilTimeWithDefault(entity.UpdatedAt, entity.CreatedAt),
		Type:      enummapper.MapAgentTypeToModel(entity.Type),
		Color:     entity.Color,
		Goal:      entity.Goal,
		IsActive:  entity.IsActive,
		Visible:   entity.VisibleInUI,
		FlowID:    utils.StringPtr(entity.FlowID),
	}
	for _, capability := range entity.CapabilitiesConfig.Capabilities {
		agentModel.Capabilities = append(agentModel.Capabilities, &model.Capability{
			ID:     capability.ID,
			Name:   capability.Name,
			Type:   enummapper.MapAgentCapabilityTypeToModel(capability.Type),
			Errors: utils.StringPtrNillable(capability.Error),
			Active: capability.Active,
			Values: capability.Values,
		})
	}
	return &agentModel
}

func MapAgentSaveInputToEntity(input model.AgentSaveInput) *postgresEntity.Agents {
	agentEntity := &postgresEntity.Agents{
		ID:          utils.IfNotNilString(input.ID),
		Name:        utils.IfNotNilString(input.Name),
		Icon:        utils.IfNotNilString(input.Icon),
		Color:       utils.IfNotNilString(input.Color),
		Goal:        utils.IfNotNilString(input.Goal),
		IsActive:    utils.IfNotNilBool(input.IsActive),
		VisibleInUI: utils.IfNotNilBool(input.Visible),
		FlowID:      utils.IfNotNilString(input.FlowID),
	}
	if input.Capabilities != nil {
		capabilities := make([]postgresEntity.Capability, 0, len(input.Capabilities))
		for _, capability := range input.Capabilities {
			capabilities = append(capabilities, postgresEntity.Capability{
				ID:     utils.IfNotNilString(capability.ID),
				Name:   utils.IfNotNilString(capability.Name),
				Error:  utils.IfNotNilString(capability.Errors),
				Active: utils.IfNotNilBool(capability.Active),
				Values: utils.IfNotNilString(capability.Values),
			})
			if capability.Type != nil {
				capabilities[len(capabilities)-1].Type = enummapper.MapAgentCapabilityTypeFromModel(*capability.Type)
			}
		}
		agentEntity.CapabilitiesConfig = postgresEntity.CapabilitiesConfig{
			Capabilities: capabilities,
		}
	}
	return agentEntity
}
