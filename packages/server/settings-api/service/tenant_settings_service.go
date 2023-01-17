package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
)

type TenantSettingsService interface {
	GetForTenant(tenantName string) (*entity.TenantSettings, error)

	SaveHubspotData(tenantName string, request dto.TenantSettingsHubspotDTO) (*entity.TenantSettings, error)
	ClearHubspotData(tenantName string) (*entity.TenantSettings, error)

	SaveZendeskData(tenantName string, request dto.TenantSettingsZendeskDTO) (*entity.TenantSettings, error)
	ClearZendeskData(tenantName string) (*entity.TenantSettings, error)

	SaveSmartSheetData(tenantName string, request dto.TenantSettingsSmartSheetDTO) (*entity.TenantSettings, error)
	ClearSmartSheetData(tenantName string) (*entity.TenantSettings, error)
}

type tenantSettingsService struct {
	repositories *repository.PostgresRepositories
}

func NewTenantSettingsService(repositories *repository.PostgresRepositories) TenantSettingsService {
	return &tenantSettingsService{
		repositories: repositories,
	}
}

func (s *tenantSettingsService) GetForTenant(tenantName string) (*entity.TenantSettings, error) {
	qr := s.repositories.TenantSettingsRepository.FindForTenantName(tenantName)
	if qr.Error != nil {
		return nil, qr.Error
	} else if qr.Result == nil {
		return nil, nil
	} else {
		settings := qr.Result.(entity.TenantSettings)
		return &settings, nil
	}
}

func (s *tenantSettingsService) SaveHubspotData(tenantName string, request dto.TenantSettingsHubspotDTO) (*entity.TenantSettings, error) {
	tenantSettings, err := s.GetForTenant(tenantName)
	if err != nil {
		return nil, err
	}

	if tenantSettings == nil {
		e := new(entity.TenantSettings)
		e.TenantName = tenantName
		e.HubspotPrivateAppKey = request.HubspotPrivateAppKey

		qr := s.repositories.TenantSettingsRepository.Save(e)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	} else {
		tenantSettings.HubspotPrivateAppKey = request.HubspotPrivateAppKey

		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	}
}

func (s *tenantSettingsService) ClearHubspotData(tenantName string) (*entity.TenantSettings, error) {
	tenantSettings, err := s.GetForTenant(tenantName)
	if err != nil {
		return nil, err
	}

	if tenantSettings == nil {
		return nil, nil
	} else {
		tenantSettings.HubspotPrivateAppKey = nil

		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	}
}

func (s *tenantSettingsService) SaveZendeskData(tenantName string, request dto.TenantSettingsZendeskDTO) (*entity.TenantSettings, error) {
	tenantSettings, err := s.GetForTenant(tenantName)
	if err != nil {
		return nil, err
	}

	if tenantSettings == nil {
		e := new(entity.TenantSettings)
		e.TenantName = tenantName
		e.ZendeskAPIKey = request.ZendeskAPIKey
		e.ZendeskAdminEmail = request.ZendeskAdminEmail
		e.ZendeskSubdomain = request.ZendeskSubdomain

		qr := s.repositories.TenantSettingsRepository.Save(e)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	} else {
		tenantSettings.ZendeskAPIKey = request.ZendeskAPIKey
		tenantSettings.ZendeskAdminEmail = request.ZendeskAdminEmail
		tenantSettings.ZendeskSubdomain = request.ZendeskSubdomain

		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	}
}

func (s *tenantSettingsService) ClearZendeskData(tenantName string) (*entity.TenantSettings, error) {
	tenantSettings, err := s.GetForTenant(tenantName)
	if err != nil {
		return nil, err
	}

	if tenantSettings == nil {
		return nil, nil
	} else {
		tenantSettings.ZendeskAPIKey = nil
		tenantSettings.ZendeskSubdomain = nil
		tenantSettings.ZendeskAdminEmail = nil

		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	}
}

func (s *tenantSettingsService) SaveSmartSheetData(tenantName string, request dto.TenantSettingsSmartSheetDTO) (*entity.TenantSettings, error) {
	tenantSettings, err := s.GetForTenant(tenantName)
	if err != nil {
		return nil, err
	}

	if tenantSettings == nil {
		e := new(entity.TenantSettings)
		e.TenantName = tenantName
		e.SmartSheetId = request.SmartSheetId
		e.SmartSheetAccessToken = request.SmartSheetAccessToken

		qr := s.repositories.TenantSettingsRepository.Save(e)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	} else {
		tenantSettings.SmartSheetId = request.SmartSheetId
		tenantSettings.SmartSheetAccessToken = request.SmartSheetAccessToken

		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	}
}

func (s *tenantSettingsService) ClearSmartSheetData(tenantName string) (*entity.TenantSettings, error) {
	tenantSettings, err := s.GetForTenant(tenantName)
	if err != nil {
		return nil, err
	}

	if tenantSettings == nil {
		return nil, nil
	} else {
		tenantSettings.SmartSheetId = nil
		tenantSettings.SmartSheetAccessToken = nil

		qr := s.repositories.TenantSettingsRepository.Save(tenantSettings)
		if qr.Error != nil {
			return nil, qr.Error
		}
		return qr.Result.(*entity.TenantSettings), nil
	}
}
