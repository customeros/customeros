package postgres_repository

import (
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type IngestEmailMessageRepository interface {
	Store(ctx context.Context, tenant, username, provider, providerMessageId string, ingestEmailMessage *postgres_entity.IngestEmailMessage) error
	UpdateState(ctx context.Context, id string, state postgres_entity.IngestEmailMessageState) error

	GetEmail(ctx context.Context, id string) (*postgres_entity.IngestEmailMessage, error)
	GetDistinctUsersForPendingMessages(ctx context.Context) ([]postgres_entity.IngestEmailMessage, error)
	GetEmailsForUserForSync(ctx context.Context, tenant, username string) ([]postgres_entity.IngestEmailMessage, error)
}

type ingestEmailMessageRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewIngestEmailMessageRepository(gormDb *gorm.DB) IngestEmailMessageRepository {
	return &ingestEmailMessageRepositoryImpl{gormDb: gormDb}
}

func (repo *ingestEmailMessageRepositoryImpl) Store(ctx context.Context, tenant, username, provider, providerMessageId string, ingestEmailMessage *postgres_entity.IngestEmailMessage) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IngestEmailMessageRepository.Store")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	tracing.TagTenant(span, tenant)

	result := postgres_entity.IngestEmailMessage{}
	err := repo.gormDb.Find(&result, "tenant = ? AND username = ? AND provider = ? AND provider_message_id = ?", tenant, username, provider, providerMessageId).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if result.Tenant != "" {
		err := errors.New("IngestEmailMessageRepository.Store - email already exists")
		tracing.TraceErr(span, err)
		return err
	}

	err = repo.gormDb.Save(&ingestEmailMessage).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (repo *ingestEmailMessageRepositoryImpl) UpdateState(ctx context.Context, id string, state postgres_entity.IngestEmailMessageState) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IngestEmailMessageRepository.UpdateState")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	err := repo.gormDb.Model(&postgres_entity.IngestEmailMessage{}).
		Where("id = ?", id).
		Update("state", state).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (repo *ingestEmailMessageRepositoryImpl) GetEmail(ctx context.Context, id string) (*postgres_entity.IngestEmailMessage, error) {
	result := postgres_entity.IngestEmailMessage{}
	err := repo.gormDb.First(&result, "id = ? ", id).Error
	if err != nil {
		tracing.TraceErr(nil, err)
		return nil, err
	}

	return &result, nil
}

func (repo *ingestEmailMessageRepositoryImpl) GetDistinctUsersForPendingMessages(ctx context.Context) ([]postgres_entity.IngestEmailMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IngestEmailMessageRepository.GetDistinctUsersForPendingMessages")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	results := []postgres_entity.IngestEmailMessage{}

	err := repo.gormDb.Select("DISTINCT tenant, username").Where("state = ?", postgres_entity.IngestEmailMessageStatePending).Find(&results).Error
	if err != nil {
		tracing.TraceErr(nil, err)
		return nil, err
	}

	return results, nil
}

func (repo *ingestEmailMessageRepositoryImpl) GetEmailsForUserForSync(ctx context.Context, tenant, username string) ([]postgres_entity.IngestEmailMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IngestEmailMessageRepository.GetEmailsForUserForSync")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)

	results := []postgres_entity.IngestEmailMessage{}

	err := repo.gormDb.Order("sent_at desc").Limit(50).Find(&results, "tenant = ? AND username = ? AND state = ?", tenant, username, postgres_entity.IngestEmailMessageStatePending).Error
	if err != nil {
		tracing.TraceErr(nil, err)
		return nil, err
	}

	return results, nil
}
