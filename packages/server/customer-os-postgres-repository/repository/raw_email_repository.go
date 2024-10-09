package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type RawEmailRepository interface {
	CountForUsername(ctx context.Context, externalSystem, tenant, username string) (int64, error)
	EmailExistsByMessageId(ctx context.Context, externalSystem, tenant, username, messageId string) (bool, error)
	Store(ctx context.Context, externalSystem, tenant, username, providerMessageId, messageId, rawEmail string, sentAt time.Time, state entity.EmailImportState) error
}

type rawEmailRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewRawEmailRepository(gormDb *gorm.DB) RawEmailRepository {
	return &rawEmailRepositoryImpl{gormDb: gormDb}
}

func (repo *rawEmailRepositoryImpl) CountForUsername(ctx context.Context, externalSystem, tenant, username string) (int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RawEmailRepository.CountForUsername")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var result int64
	err := repo.gormDb.Model(entity.RawEmail{}).Where("external_system = ? AND tenant = ? AND username = ?", externalSystem, tenant, username).Count(&result).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return 0, err
	}

	return result, nil
}

func (repo *rawEmailRepositoryImpl) EmailExistsByMessageId(ctx context.Context, externalSystem, tenant, username, messageId string) (bool, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RawEmailRepository.EmailExistsByMessageId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var result int64
	err := repo.gormDb.Model(entity.RawEmail{}).Where("external_system = ? AND tenant = ? AND username = ? AND message_id = ?", externalSystem, tenant, username, messageId).Count(&result).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	return result > 0, nil
}

func (repo *rawEmailRepositoryImpl) Store(ctx context.Context, externalSystem, tenant, username, providerMessageId, messageId, rawEmail string, sentAt time.Time, state entity.EmailImportState) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RawEmailRepository.Store")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	result := entity.RawEmail{}
	err := repo.gormDb.Find(&result, "external_system = ? AND tenant = ? AND username = ? AND message_id = ?", externalSystem, tenant, username, messageId).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if result.Tenant != "" {
		err := errors.New("RawEmailRepository.Store - email already exists")
		tracing.TraceErr(span, err)
		return err
	}

	result.ProviderMessageId = providerMessageId
	result.MessageId = messageId

	result.CreatedAt = utils.Now()
	result.SentAt = sentAt
	result.State = state
	result.ExternalSystem = externalSystem
	result.Tenant = tenant
	result.Username = username
	result.Data = rawEmail
	result.SentToEventStoreState = "PENDING"

	err = repo.gormDb.Save(&result).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
