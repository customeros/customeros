package repository

import (
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type MailstackBuyRequestRepository interface {
	GetList(ctx context.Context) ([]*entity.MailstackBuyRequest, error)
	GetById(ctx context.Context, id string) (*entity.MailstackBuyRequest, error)
	Store(ctx context.Context, tx *gorm.DB, input *entity.MailstackBuyRequest) (string, error)

	GetDomains(ctx context.Context, mailstackBuyRequestId string) ([]*entity.MailstackBuyRequestDomain, error)
	StoreDomain(ctx context.Context, tx *gorm.DB, input *entity.MailstackBuyRequestDomain) error

	GetMailboxes(ctx context.Context, mailstackBuyRequestId string) ([]*entity.MailstackBuyRequestMailbox, error)
	StoreMailbox(ctx context.Context, tx *gorm.DB, input *entity.MailstackBuyRequestMailbox) error
}

type mailstackBuyRequestRepositoryRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewMailstackBuyRequestRepository(gormDb *gorm.DB) MailstackBuyRequestRepository {
	return &mailstackBuyRequestRepositoryRepositoryImpl{gormDb: gormDb}
}

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) GetList(ctx context.Context) ([]*entity.MailstackBuyRequest, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackBuyRequestRepository.GetList")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)

	var e []*entity.MailstackBuyRequest
	err := repo.gormDb.Where("tenant = ?", tenant).Find(&e).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return e, nil
}

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) GetById(ctx context.Context, id string) (*entity.MailstackBuyRequest, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackBuyRequestRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	var e *entity.MailstackBuyRequest
	err := repo.gormDb.Where("tenant = ? and id = ?", tenant, id).First(&e).Error
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

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) Store(ctx context.Context, tx *gorm.DB, input *entity.MailstackBuyRequest) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackBuyRequestRepository.Store")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(log.String("input.domains", input.Domains), log.String("input.usernames", input.Usernames))

	if input.Domains == "" || input.Usernames == "" {
		span.LogFields(log.Object("input", input))
		err := errors.New("params missing")
		tracing.TraceErr(span, err)
		return "", err
	}

	input.Tenant = tenant
	input.CreatedAt = utils.Now()

	if tx == nil {
		tx = repo.gormDb
	}

	err := repo.gormDb.Save(&input).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return input.ID, nil
}

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) GetDomains(ctx context.Context, mailstackBuyRequestId string) ([]*entity.MailstackBuyRequestDomain, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackBuyRequestRepository.GetDomains")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	span.LogFields(log.String("mailstackBuyRequestId", mailstackBuyRequestId))

	tenant := common.GetTenantFromContext(ctx)

	var e []*entity.MailstackBuyRequestDomain
	err := repo.gormDb.Where("tenant = ? and mailstack_buy_request_id = ?", tenant, mailstackBuyRequestId).Find(&e).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return e, nil
}

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) StoreDomain(ctx context.Context, tx *gorm.DB, input *entity.MailstackBuyRequestDomain) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackBuyRequestRepository.StoreDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(log.String("input.domain", input.Domain))

	if input.Domain == "" || input.MailstackBuyRequestId == "" {
		span.LogFields(log.Object("input", input))
		err := errors.New("params missing")
		tracing.TraceErr(span, err)
		return err
	}

	input.Tenant = tenant
	input.CreatedAt = utils.Now()

	if tx == nil {
		tx = repo.gormDb
	}

	err := repo.gormDb.Save(&input).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) GetMailboxes(ctx context.Context, mailstackBuyRequestId string) ([]*entity.MailstackBuyRequestMailbox, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackBuyRequestRepository.GetMailboxes")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	span.LogFields(log.String("mailstackBuyRequestId", mailstackBuyRequestId))

	tenant := common.GetTenantFromContext(ctx)

	var e []*entity.MailstackBuyRequestMailbox
	err := repo.gormDb.Where("tenant = ? and mailstack_buy_request_id = ?", tenant, mailstackBuyRequestId).Find(&e).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return e, nil
}

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) StoreMailbox(ctx context.Context, tx *gorm.DB, input *entity.MailstackBuyRequestMailbox) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailstackBuyRequestRepository.StoreMailbox")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)

	span.LogFields(log.String("input.domain", input.Domain), log.String("input.username", input.Username), log.String("input.mailbox", input.Mailbox), log.String("input.mailstackBuyRequestId", input.MailstackBuyRequestId))

	if input.Domain == "" || input.Username == "" || input.Mailbox == "" || input.MailstackBuyRequestId == "" {
		span.LogFields(log.Object("input", input))
		err := errors.New("params missing")
		tracing.TraceErr(span, err)
		return err
	}

	input.Tenant = tenant
	input.CreatedAt = utils.Now()

	if tx == nil {
		tx = repo.gormDb
	}

	err := repo.gormDb.Save(&input).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
