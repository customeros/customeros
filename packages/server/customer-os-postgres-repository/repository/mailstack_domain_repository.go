package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type MailStackDomainRepository interface {
	RegisterDomain(ctx context.Context, tenant, domain string) (*entity.MailStackDomain, error)
}

type mailStackDomainRepository struct {
	db *gorm.DB
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
	}

	err := r.db.Create(&mailStackDomain).Error
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return nil, err
	}

	return &mailStackDomain, nil
}
