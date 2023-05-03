package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/entity"
	"gorm.io/gorm"
)

const GMAIL_SERVICE_PRIVATE_KEY = "GMAIL_SERVICE_PRIVATE_KEY"
const GMAIL_SERVICE_EMAIL_ADDRESS = "GMAIL_SERVICE_EMAIL_ADDRESS"

type ApiKeyRepository interface {
	GetApiKeyByTenantService(tenantId, serviceId string) (string, error)
}

type ApiKeyRepositoryImpl struct {
	StorageDB *config.StorageDB
}

func NewApiKeyRepository(db *config.StorageDB) *ApiKeyRepositoryImpl {
	return &ApiKeyRepositoryImpl{StorageDB: db}
}

func (repo *ApiKeyRepositoryImpl) GetApiKeyByTenantService(tenantName, serviceId string) (string, error) {
	result := entity.TenantAPIKey{}
	err := repo.StorageDB.GormDB.First(&result, "tenant_name = ? AND key = ?", tenantName, serviceId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// handle record not found error
			return "", nil
		} else {
			return "", fmt.Errorf("GetApiKeyByTenantService: %s", err.Error())
		}
	}
	return result.Value, nil
}
