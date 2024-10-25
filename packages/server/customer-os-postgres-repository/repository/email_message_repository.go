package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type EmailMessageRepository interface {
	GetByProviderMessageId(ctx context.Context, tenant, providerMessageId string) (*entity.EmailMessage, error)
	GetByProducer(ctx context.Context, tenant, producerId, producerType string) (*entity.EmailMessage, error)
	GetForSending(ctx context.Context) ([]*entity.EmailMessage, error)
	GetForProcessing(ctx context.Context) ([]*entity.EmailMessage, error)
	Store(ctx context.Context, tenant string, input *entity.EmailMessage) error
}

type emailMessageRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewEmailMessageRepository(gormDb *gorm.DB) EmailMessageRepository {
	return &emailMessageRepositoryImpl{gormDb: gormDb}
}

func (repo *emailMessageRepositoryImpl) GetByProviderMessageId(ctx context.Context, tenant, providerMessageId string) (*entity.EmailMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailMessageRepository.GetByProviderMessageId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	span.LogFields(log.String("tenant", tenant), log.String("providerMessageId", providerMessageId))

	if tenant == "" || providerMessageId == "" {
		err := errors.New("params missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var e *entity.EmailMessage
	err := repo.gormDb.Where("tenant = ? and provider_message_id = ?", tenant, providerMessageId).First(&e).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			span.LogFields(log.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return e, nil
}

func (repo *emailMessageRepositoryImpl) GetByProducer(ctx context.Context, tenant, producerId, producerType string) (*entity.EmailMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailMessageRepository.GetByProducer")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	span.LogFields(log.String("tenant", tenant), log.String("producerId", producerId), log.String("producerType", producerType))

	if tenant == "" || producerId == "" || producerType == "" {
		err := errors.New("params missing")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var e *entity.EmailMessage
	err := repo.gormDb.Where("tenant = ? and producer_id = ? and producer_type = ?", tenant, producerId, producerType).First(&e).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			span.LogFields(log.Bool("result.found", false))
			return nil, nil
		}
		tracing.TraceErr(span, err)
		return nil, err
	}

	return e, nil
}

func (repo *emailMessageRepositoryImpl) GetForSending(ctx context.Context) ([]*entity.EmailMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailMessageRepository.GetForSending")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var entities []*entity.EmailMessage
	err := repo.gormDb.Where("status = ?", entity.EmailMessageStatusScheduled).Order("created_at asc").Limit(25).Find(&entities).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return entities, nil
}

func (repo *emailMessageRepositoryImpl) GetForProcessing(ctx context.Context) ([]*entity.EmailMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailMessageRepository.GetForProcessing")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var entities []*entity.EmailMessage
	err := repo.gormDb.Where("status = ?", entity.EmailMessageStatusSent).Order("created_at asc").Limit(25).Find(&entities).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return entities, nil
}

func (repo *emailMessageRepositoryImpl) Store(ctx context.Context, tenant string, input *entity.EmailMessage) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailMessageRepository.Store")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	span.LogFields(log.String("tenant", tenant), log.String("producerId", input.ProducerId), log.String("producerType", input.ProducerType))

	if input.Status == "" || input.ProducerId == "" || input.ProducerType == "" || input.From == "" || len(input.To) == 0 || input.Subject == "" || input.Content == "" {
		err := errors.New("params missing")
		tracing.TraceErr(span, err)
		return err
	}

	input.Tenant = tenant
	input.CreatedAt = utils.Now()

	err := repo.gormDb.Save(&input).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
