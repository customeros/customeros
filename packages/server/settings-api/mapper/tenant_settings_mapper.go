package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
)

func MapTenantSettingsEntityToDTO(input *entity.TenantSettings) *dto.TenantSettingsResponseDTO {
	responseDTO := dto.TenantSettingsResponseDTO{}

	if input != nil && input.HubspotPrivateAppKey != nil {
		responseDTO.HubspotExists = true
	}

	if input != nil && input.ZendeskAPIKey != nil && input.ZendeskSubdomain != nil && input.ZendeskAdminEmail != nil {
		responseDTO.ZendeskExists = true
	}

	if input != nil && input.SmartSheetId != nil && input.SmartSheetAccessToken != nil {
		responseDTO.SmartSheetExists = true
	}

	return &responseDTO
}
