package postgres_repository

import (
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type EmailMessageRepository interface {
	GetByProviderMessageId(ctx context.Context, tenant, providerMessageId string) (*postgres_entity.EmailMessage, error)
	GetByProducer(ctx context.Context, tenant, producerId, producerType string) (*postgres_entity.EmailMessage, error)
	GetForSending(ctx context.Context) ([]*postgres_entity.EmailMessage, error)
	GetForProcessing(ctx context.Context) ([]*postgres_entity.EmailMessage, error)
	Store(ctx context.Context, tenant string, input *postgres_entity.EmailMessage) error
}

type emailMessageRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewEmailMessageRepository(gormDb *gorm.DB) EmailMessageRepository {
	return &emailMessageRepositoryImpl{gormDb: gormDb}
}

func (repo *emailMessageRepositoryImpl) GetByProviderMessageId(ctx context.Context, tenant, providerMessageId string) (*postgres_entity.EmailMessage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "EmailMessageRepository.GetByProviderMessageId")
	defer spans.Finish()

	spans.LogKV("tenant", tenant)
	spans.LogKV("providerMessageId", providerMessageId)

	if tenant == "" || providerMessageId == "" {
		err := errors.New("params missing")
		spans.TraceError(err)
		return nil, err
	}

	var e *postgres_entity.EmailMessage
	err := repo.gormDb.Where("tenant = ? and provider_message_id = ?", tenant, providerMessageId).First(&e).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	return e, nil
}

func (repo *emailMessageRepositoryImpl) GetByProducer(ctx context.Context, tenant, producerId, producerType string) (*postgres_entity.EmailMessage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "EmailMessageRepository.GetByProducer")
	defer spans.Finish()

	spans.LogKV("tenant", tenant)
	spans.LogKV("producerId", producerId)
	spans.LogKV("producerType", producerType)

	if tenant == "" || producerId == "" || producerType == "" {
		err := errors.New("params missing")
		spans.TraceError(err)
		return nil, err
	}

	var e *postgres_entity.EmailMessage
	err := repo.gormDb.Where("tenant = ? and producer_id = ? and producer_type = ?", tenant, producerId, producerType).First(&e).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			spans.LogKV("result.found", false)
			return nil, nil
		}
		spans.TraceError(err)
		return nil, err
	}

	return e, nil
}

func (repo *emailMessageRepositoryImpl) GetForSending(ctx context.Context) ([]*postgres_entity.EmailMessage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "EmailMessageRepository.GetForSending")
	defer spans.Finish()

	var entities []*postgres_entity.EmailMessage
	err := repo.gormDb.Where("status = ?", postgres_entity.EmailMessageStatusScheduled).Order("created_at asc").Limit(25).Find(&entities).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return entities, nil
}

func (repo *emailMessageRepositoryImpl) GetForProcessing(ctx context.Context) ([]*postgres_entity.EmailMessage, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "EmailMessageRepository.GetForProcessing")
	defer spans.Finish()

	var entities []*postgres_entity.EmailMessage
	err := repo.gormDb.Where("status = ?", postgres_entity.EmailMessageStatusSent).Order("created_at asc").Limit(25).Find(&entities).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return entities, nil
}

func (repo *emailMessageRepositoryImpl) Store(ctx context.Context, tenant string, input *postgres_entity.EmailMessage) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "EmailMessageRepository.Store")
	defer spans.Finish()

	spans.LogKV("tenant", tenant)
	spans.LogKV("producerId", input.ProducerId)
	spans.LogKV("producerType", input.ProducerType)

	if input.Status == "" || input.ProducerId == "" || input.ProducerType == "" || input.From == "" || len(input.To) == 0 || input.Subject == "" || input.Content == "" {
		spans.LogObjectAsJson("input", input)
		err := errors.New("params missing")
		spans.TraceError(err)
		return err
	}

	input.Tenant = tenant
	input.CreatedAt = utils.Now()

	err := repo.gormDb.Save(&input).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
