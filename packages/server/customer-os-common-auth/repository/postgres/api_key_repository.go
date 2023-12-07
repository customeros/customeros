package repository

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type ApiKeyRepository interface {
	GetApiKeyByTenantService(ctx context.Context, tenantId, serviceId string) (string, error)
}

type ApiKeyRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewApiKeyRepository(gormDb *gorm.DB) *ApiKeyRepositoryImpl {
	return &ApiKeyRepositoryImpl{gormDb: gormDb}
}

func (repo *ApiKeyRepositoryImpl) GetApiKeyByTenantService(ctx context.Context, tenantName, serviceId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ApiKeyRepository.GetApiKeyByTenantService")
	defer span.Finish()
	span.LogFields(log.String("tenantName", tenantName), log.String("serviceId", serviceId))

	result := entity.TenantAPIKey{}
	err := repo.gormDb.First(&result, "tenant_name = ? AND key = ?", tenantName, serviceId).Error

	if err != nil {
		tracing.TraceErr(span, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.LogFields(log.String("result", "record not found"))
			// handle record not found error
			return "", nil
		} else {
			return "", fmt.Errorf("GetApiKeyByTenantService: %s", err.Error())
		}
	}
	return result.Value, nil
}
