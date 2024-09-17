package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type MailStackDomainRepository interface {
	RegisterDomain(ctx context.Context, tenant, domain string) (*entity.MailStackDomain, error)
	CheckDomainOwnership(ctx context.Context, tenant, domain string) (bool, error)
	GetActiveDomains(ctx context.Context, tenant string) ([]entity.MailStackDomain, error)
	SetConfigured(ctx context.Context, tenant, domain string) error
}

type mailStackDomainRepository struct {
	db *gorm.DB
}

func (r *mailStackDomainRepository) SetConfigured(ctx context.Context, tenant, domain string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "MailStackDomainRepository.SetConfigured")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("domain", domain)

	err := r.db.WithContext(ctx).
		Model(&entity.MailStackDomain{}).
		Where("tenant = ? AND domain = ?", tenant, domain).
		UpdateColumn("configured", true).
		UpdateColumn("updated_at", utils.Now()).
		Error
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return err
	}

	return nil
}

func NewMailStackDomainRepository(db *gorm.DB) MailStackDomainRepository {
	return &mailStackDomainRepository{db: db}
}

func (r *mailStackDomainRepository) RegisterDomain(ctx context.Context, tenant, domain string) (*entity.MailStackDomain, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "MailStackDomainRepository.RegisterDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	now := utils.Now()
	mailStackDomain := entity.MailStackDomain{
		Tenant:    tenant,
		Domain:    domain,
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}

	err := r.db.Create(&mailStackDomain).Error
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return nil, err
	}

	return &mailStackDomain, nil
}

func (r *mailStackDomainRepository) CheckDomainOwnership(ctx context.Context, tenant, domain string) (bool, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "MailStackDomainRepository.CheckDomainOwnership")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("domain", domain)

	var mailStackDomain entity.MailStackDomain
	err := r.db.WithContext(ctx).
		Where("tenant = ? AND domain = ? AND active = ?", tenant, domain, true).
		First(&mailStackDomain).Error

	if err != nil {
		// If the record is not found, return false without an error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.LogFields(tracingLog.Bool("response.exists", false))
			return false, nil
		}
		// If any other error occurs, log and trace it
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return false, err
	}

	// If the record is found, return true
	span.LogFields(tracingLog.Bool("response.exists", true))
	return true, nil
}

func (r *mailStackDomainRepository) GetActiveDomains(ctx context.Context, tenant string) ([]entity.MailStackDomain, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "MailStackDomainRepository.GetActiveDomains")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	var mailStackDomains []entity.MailStackDomain
	err := r.db.WithContext(ctx).
		Where("tenant = ? AND active = ?", tenant, true).
		Find(&mailStackDomains).Error
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return nil, err
	}

	return mailStackDomains, nil
}
