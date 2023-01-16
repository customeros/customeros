package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
)

func MapTenantSettingsEntityToDTO(input *entity.TenantSettings) *dto.TenantSettingsDTO {
	if input == nil {
		return nil
	}
	file := dto.TenantSettingsDTO{
		Id:                   input.ID,
		HubspotPrivateAppKey: input.HubspotPrivateAppKey,
		ZendeskAPIKey:        input.ZendeskAPIKey,
		ZendeskSubdomain:     input.ZendeskSubdomain,
		ZendeskAdminEmail:    input.ZendeskAdminEmail,
	}
	return &file
}

func MapTenantSettingsDTOToEntity(dto *dto.TenantSettingsDTO, tenantId string) *entity.TenantSettings {
	if dto == nil {
		return nil
	}
	tenantSettings := entity.TenantSettings{
		ID:                   dto.Id,
		TenantId:             tenantId,
		HubspotPrivateAppKey: dto.HubspotPrivateAppKey,
		ZendeskAPIKey:        dto.ZendeskAPIKey,
		ZendeskSubdomain:     dto.ZendeskSubdomain,
		ZendeskAdminEmail:    dto.ZendeskAdminEmail,
	}
	return &tenantSettings
}
