package postgres_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type MailStackDomainRepository interface {
	CreateDMARCReport(ctx context.Context, tenant string, report *postgres_entity.DMARCMonitoring) error
	GetDomainCrossTenant(ctx context.Context, domain string) (*postgres_entity.MailStackDomain, error)
}

type mailStackDomainRepository struct {
	db *gorm.DB
}

func NewMailStackDomainRepository(db *gorm.DB) MailStackDomainRepository {
	return &mailStackDomainRepository{db: db}
}

func (r *mailStackDomainRepository) CreateDMARCReport(ctx context.Context, tenant string, report *postgres_entity.DMARCMonitoring) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "MailStackDomainRepository.CreateDMARCReport")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	now := utils.Now()
	report.CreatedAt = now

	err := r.db.Create(&report).Error
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return err
	}
	return nil
}

func (r *mailStackDomainRepository) GetDomainCrossTenant(ctx context.Context, domain string) (*postgres_entity.MailStackDomain, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "MailStackDomainRepository.GetDomainCrossTenant")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("domain", domain)

	var mailStackDomain postgres_entity.MailStackDomain
	err := r.db.WithContext(ctx).
		Where("domain = ?", domain).
		First(&mailStackDomain).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return nil, err
	}

	return &mailStackDomain, nil
}
