package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.LogFields(log.String("result", "record not found"))
			return "", nil
		} else {
			tracing.TraceErr(span, err)
			return "", errors.Wrap(err, "GetApiKeyByTenantService")
		}
	}
	return result.Value, nil
}
