package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type GoogleServiceAccountKeyRepository interface {
	GetApiKeyByTenantService(ctx context.Context, tenantId, serviceId string) (string, error)
	SaveKey(ctx context.Context, tenant, key, value string) error
	DeleteKey(ctx context.Context, tenant, key string) error
}

type GoogleServiceAccountKeyRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewGoogleServiceAccountKeyRepository(gormDb *gorm.DB) *GoogleServiceAccountKeyRepositoryImpl {
	return &GoogleServiceAccountKeyRepositoryImpl{gormDb: gormDb}
}

func (repo *GoogleServiceAccountKeyRepositoryImpl) GetApiKeyByTenantService(ctx context.Context, tenantName, serviceId string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GoogleServiceAccountKeyRepository.GetApiKeyByTenantService")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(log.String("tenantName", tenantName), log.String("serviceId", serviceId))

	result := entity.GoogleServiceAccountKey{}
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

func (repo *GoogleServiceAccountKeyRepositoryImpl) SaveKey(ctx context.Context, tenant, key, value string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleServiceAccountKeyRepository.SaveKey")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("key", key), log.String("value", value))

	existing, err := repo.GetApiKeyByTenantService(ctx, tenant, key)
	if err != nil {
		return err
	}

	if existing != "" {
		return nil
	}

	newKey := entity.GoogleServiceAccountKey{
		TenantName: tenant,
		Key:        key,
		Value:      value,
	}

	result := repo.gormDb.Save(&newKey)

	if result.Error != nil {
		return errors.Wrap(result.Error, "SaveKey")
	}

	return nil
}

func (repo *GoogleServiceAccountKeyRepositoryImpl) DeleteKey(ctx context.Context, tenant, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleServiceAccountKeyRepository.SaveKey")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogFields(log.String("key", key))

	existing, err := repo.GetApiKeyByTenantService(ctx, tenant, key)
	if err != nil {
		return err
	}

	if existing != "" {
		return nil
	}

	err = repo.gormDb.Delete(&entity.GoogleServiceAccountKey{}, "tenant_name = ? AND key = ?", tenant, key).Error
	if err != nil {
		return errors.Wrap(err, "DeleteKey")
	}

	return nil
}
