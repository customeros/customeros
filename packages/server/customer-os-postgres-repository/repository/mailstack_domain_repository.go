package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type MailStackDomainRepository interface {
	RegisterDomain(ctx context.Context, tenant, domain string, createdAt *time.Time) (*entity.MailStackDomain, error)
}

type mailStackDomainRepository struct {
	db *gorm.DB
}

func NewMailStackDomainRepository(db *gorm.DB) MailStackDomainRepository {
	return &mailStackDomainRepository{db: db}
}

func (r *mailStackDomainRepository) RegisterDomain(ctx context.Context, tenant, domain string, createdAt *time.Time) (*entity.MailStackDomain, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "MailStackDomainRepository.RegisterDomain")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.TagTenant(span, tenant)

	mailStackDomain := entity.MailStackDomain{
		Tenant: tenant,
		Domain: domain,
	}
	if createdAt != nil {
		mailStackDomain.CreatedAt = *createdAt
		mailStackDomain.UpdatedAt = *createdAt
	} else {
		now := utils.Now()
		mailStackDomain.CreatedAt = now
		mailStackDomain.UpdatedAt = now
	}

	err := r.db.Create(&mailStackDomain).Error
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "db error"))
		return nil, err
	}

	return &mailStackDomain, nil
}
