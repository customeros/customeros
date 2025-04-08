package postgres_repository

import (
	"errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type MailstackBuyRequestRepository interface {
	GetList(ctx context.Context) ([]*postgres_entity.MailstackBuyRequest, error)
	GetById(ctx context.Context, id string) (*postgres_entity.MailstackBuyRequest, error)
	Store(ctx context.Context, tx *gorm.DB, input *postgres_entity.MailstackBuyRequest) (string, error)

	GetDomains(ctx context.Context, mailstackBuyRequestId string) ([]*postgres_entity.MailstackBuyRequestDomain, error)
	StoreDomain(ctx context.Context, tx *gorm.DB, input *postgres_entity.MailstackBuyRequestDomain) error
}

type mailstackBuyRequestRepositoryRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewMailstackBuyRequestRepository(gormDb *gorm.DB) MailstackBuyRequestRepository {
	return &mailstackBuyRequestRepositoryRepositoryImpl{gormDb: gormDb}
}

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) GetList(ctx context.Context) ([]*postgres_entity.MailstackBuyRequest, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "MailstackBuyRequestRepository.GetList")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	var e []*postgres_entity.MailstackBuyRequest
	err := repo.gormDb.Where("tenant = ?", tenant).Find(&e).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return e, nil
}

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) GetById(ctx context.Context, id string) (*postgres_entity.MailstackBuyRequest, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "MailstackBuyRequestRepository.GetById")
	defer spans.Finish()
	spans.LogKV("id", id)

	tenant := common.GetTenantFromContext(ctx)

	var e *postgres_entity.MailstackBuyRequest
	err := repo.gormDb.Where("tenant = ? and id = ?", tenant, id).First(&e).Error
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

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) Store(ctx context.Context, tx *gorm.DB, input *postgres_entity.MailstackBuyRequest) (string, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "MailstackBuyRequestRepository.Store")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	spans.LogObjectAsJson("input", input)

	if input.Domains == "" || input.Usernames == "" {
		err := errors.New("params missing")
		spans.TraceError(err)
		return "", err
	}

	input.Tenant = tenant
	input.CreatedAt = utils.Now()

	if tx == nil {
		tx = repo.gormDb
	}

	err := repo.gormDb.Save(&input).Error
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	return input.ID, nil
}

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) GetDomains(ctx context.Context, mailstackBuyRequestId string) ([]*postgres_entity.MailstackBuyRequestDomain, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "MailstackBuyRequestRepository.GetDomains")
	defer spans.Finish()
	spans.LogKV("mailstackBuyRequestId", mailstackBuyRequestId)

	tenant := common.GetTenantFromContext(ctx)

	var e []*postgres_entity.MailstackBuyRequestDomain
	err := repo.gormDb.Where("tenant = ? and mailstack_buy_request_id = ?", tenant, mailstackBuyRequestId).Find(&e).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return e, nil
}

func (repo *mailstackBuyRequestRepositoryRepositoryImpl) StoreDomain(ctx context.Context, tx *gorm.DB, input *postgres_entity.MailstackBuyRequestDomain) error {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "MailstackBuyRequestRepository.StoreDomain")
	defer spans.Finish()
	spans.LogObjectAsJson("input", input)

	tenant := common.GetTenantFromContext(ctx)

	if input.Domain == "" || input.MailstackBuyRequestId == "" {
		err := errors.New("params missing")
		spans.TraceError(err)
		return err
	}

	input.Tenant = tenant
	input.CreatedAt = utils.Now()

	if tx == nil {
		tx = repo.gormDb
	}

	err := repo.gormDb.Save(&input).Error
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
