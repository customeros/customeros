package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AiPromptLogRepository interface {
	Store(aiPrompt entity.AiPromptLog) (string, error)
	UpdateResponse(id string, rawResponse string) error
	UpdateError(id string, postProcessErrorMessage string) error
}

type aiPromptLogRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewAiPromptLogRepository(gormDb *gorm.DB) AiPromptLogRepository {
	return &aiPromptLogRepositoryImpl{gormDb: gormDb}
}

func (repo *aiPromptLogRepositoryImpl) Store(aiPrompt entity.AiPromptLog) (string, error) {
	err := repo.gormDb.Save(&aiPrompt).Error

	if err != nil {
		logrus.Errorf("Failed storing aiPromptResponse: %v", aiPrompt)
		return "", err
	}

	return aiPrompt.ID, nil
}

func (repo *aiPromptLogRepositoryImpl) UpdateResponse(id string, rawResponse string) error {
	err := repo.gormDb.Model(&entity.AiPromptLog{}).Where("id = ?", id).Update("raw_response", rawResponse).Error

	if err != nil {
		logrus.Errorf("Failed marking email as sent to event store: %v", id)
		return err
	}

	return nil
}

func (repo *aiPromptLogRepositoryImpl) UpdateError(id string, postProcessErrorMessage string) error {
	tx := repo.gormDb.Model(&entity.AiPromptLog{}).Where("id = ?", id)

	if postProcessErrorMessage != "" {
		tx.Update("post_process_error", true)
		tx.Update("post_process_error_message", postProcessErrorMessage)
	} else {
		tx.Update("post_process_error", false)
		tx.Update("post_process_error_message", nil)
	}

	err := tx.Error

	if err != nil {
		logrus.Errorf("Failed marking email as sent to event store: %v", id)
		return err
	}

	return nil
}
