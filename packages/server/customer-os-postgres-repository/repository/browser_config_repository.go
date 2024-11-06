package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type BrowserConfigRepository interface {
	Get(ctx context.Context) ([]entity.BrowserConfig, error)
	GetForUser(ctx context.Context, userId string) (*entity.BrowserConfig, error)

	Merge(ctx context.Context, browserConfig *entity.BrowserConfig) error
}

type browserConfigRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewBrowserConfigRepository(gormDb *gorm.DB) BrowserConfigRepository {
	return &browserConfigRepositoryImpl{gormDb: gormDb}
}

func (repo *browserConfigRepositoryImpl) Get(ctx context.Context) ([]entity.BrowserConfig, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserConfigRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result []entity.BrowserConfig
	err := repo.gormDb.Where("session_status = 'VALID'").Find(&result).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}

func (repo *browserConfigRepositoryImpl) GetForUser(ctx context.Context, userId string) (*entity.BrowserConfig, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserConfigRepository.Get")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)

	var result *entity.BrowserConfig
	err := repo.gormDb.Where("tenant = ? and user_id = ? and session_status = 'VALID'", tenant, userId).First(&result).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}

func (r *browserConfigRepositoryImpl) Merge(ctx context.Context, input *entity.BrowserConfig) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "BrowserConfigRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(tracingLog.Object("input", input))

	// Check if the browserConfig already exists
	var browserConfig entity.BrowserConfig
	err := r.gormDb.
		Where("tenant = ? AND user_id = ?", tenant, input.UserId).
		First(&browserConfig).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tracing.TraceErr(span, err)
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		browserConfig = entity.BrowserConfig{
			Tenant: tenant,
			UserId: input.UserId,
			Status: input.Status,
		}

		err = r.gormDb.Create(&browserConfig).Error
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		browserConfig.UserId = input.UserId
		browserConfig.Status = input.Status

		err = r.gormDb.Save(&browserConfig).Error
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
