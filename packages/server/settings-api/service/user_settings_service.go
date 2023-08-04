package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository"
)

type UserSettingsService interface {
	GetByUserName(userName string) (*model.UserSettings, error)
	Save(settings *model.UserSettings)
	Delete(tenantName, identifier string) error
}

type userSettingsService struct {
	repositories *repository.PostgresRepositories
	log          logger.Logger
}

func NewUserSettingsService(repositories *repository.PostgresRepositories, log logger.Logger) UserSettingsService {
	return &userSettingsService{
		repositories: repositories,
		log:          log,
	}
}

func (u userSettingsService) GetByUserName(userName string) (*model.UserSettings, error) {
	qr := u.repositories.UserSettingRepository.GetByUserName(userName)
	var settings entity.UserSettingsEntity
	var ok bool
	if qr.Error != nil {
		return nil, qr.Error
	} else if qr.Result == nil {
		return nil, nil
	} else {
		settings, ok = qr.Result.(entity.UserSettingsEntity)
		if !ok {
			return nil, fmt.Errorf("GetForTenant: unexpected type %T", qr.Result)
		}
	}

	return mapUserSettingsEntityToDTO(&settings), nil
}

func (u userSettingsService) Save(userSettings *model.UserSettings) {
	u.repositories.UserSettingRepository.Save(mapUserSettingsDTOToEntity(userSettings))
}

func (u userSettingsService) Delete(tenantName, identifier string) error {
	//TODO implement me
	panic("implement me")
}

func mapUserSettingsDTOToEntity(dto *model.UserSettings) *entity.UserSettingsEntity {
	return &entity.UserSettingsEntity{
		ID:                          uuid.MustParse(dto.ID),
		TenantName:                  dto.TenantName,
		UserName:                    dto.UserName,
		GoogleOAuthAllScopesEnabled: dto.GoogleOAuthAllScopesEnabled,
		GoogleOAuthUserAccessToken:  dto.GoogleOAuthUserAccessToken,
	}
}

func mapUserSettingsEntityToDTO(entity *entity.UserSettingsEntity) *model.UserSettings {
	return &model.UserSettings{
		ID:                          entity.ID.String(),
		TenantName:                  entity.TenantName,
		UserName:                    entity.UserName,
		GoogleOAuthAllScopesEnabled: entity.GoogleOAuthAllScopesEnabled,
		GoogleOAuthUserAccessToken:  entity.GoogleOAuthUserAccessToken,
	}
}
