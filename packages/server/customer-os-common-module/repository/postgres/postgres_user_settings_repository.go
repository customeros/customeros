package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/repository/helper"
	"gorm.io/gorm"
)

type UserSettingsRepository interface {
	GetByUserName(username string) helper.QueryResult
	Save(userSettings *entity.UserSettings) helper.QueryResult
	Delete(username string) error
}

type userSettingsRepo struct {
	db *gorm.DB
}

func NewUserSettingsRepository(db *gorm.DB) UserSettingsRepository {
	return &userSettingsRepo{
		db: db,
	}
}

func (u userSettingsRepo) GetByUserName(username string) helper.QueryResult {
	//TODO implement me
	panic("implement me")
}

func (u userSettingsRepo) Save(userSettings *entity.UserSettings) helper.QueryResult {
	//TODO implement me
	panic("implement me")
}

func (u userSettingsRepo) Delete(username string) error {
	//TODO implement me
	panic("implement me")
}
