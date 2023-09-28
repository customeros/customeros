package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"gorm.io/gorm"
)

type ApiKeyRepository interface {
	GetApiKeyByTenantService(tenantId, serviceId string) (string, error)
}

type ApiKeyRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewApiKeyRepository(gormDb *gorm.DB) *ApiKeyRepositoryImpl {
	return &ApiKeyRepositoryImpl{gormDb: gormDb}
}

func (repo *ApiKeyRepositoryImpl) GetApiKeyByTenantService(tenantName, serviceId string) (string, error) {
	result := entity.TenantAPIKey{}
	err := repo.gormDb.First(&result, "tenant_name = ? AND key = ?", tenantName, serviceId).Error

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
