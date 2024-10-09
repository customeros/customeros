package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type EmailMessageRepository interface {
	GetForSending(ctx context.Context) ([]*entity.EmailMessage, error)
	Store(ctx context.Context, tenant string, input *entity.EmailMessage) error
	DeleteByProducerId(ctx context.Context, tenant, producerId string) error
}

type emailMessageRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewEmailMessageRepository(gormDb *gorm.DB) EmailMessageRepository {
	return &emailMessageRepositoryImpl{gormDb: gormDb}
}

func (repo *emailMessageRepositoryImpl) GetForSending(ctx context.Context) ([]*entity.EmailMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailMessageRepository.GetForSending")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var entities []*entity.EmailMessage
	err := repo.gormDb.Find(&entities).Where("sent_at is null and error is null").Order("created_at asc").Limit(25).Error
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

	if input.ProducerId == "" || input.ProducerType == "" || input.From == "" || len(input.To) == 0 || input.Subject == "" || input.Content == "" {
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

func (repo *emailMessageRepositoryImpl) DeleteByProducerId(ctx context.Context, tenant, producerId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailMessageRepository.DeleteByProducerId")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	if tenant == "" || producerId == "" {
		err := errors.New("params missing")
		tracing.TraceErr(span, err)
		return err
	}

	err := repo.gormDb.Delete(&entity.EmailMessage{}, "tenant = ? AND producer_id = ?", tenant, producerId).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
