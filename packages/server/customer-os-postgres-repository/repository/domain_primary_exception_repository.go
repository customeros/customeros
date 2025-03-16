package postgres_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

type DomainPrimaryExceptionRepository interface {
	Exists(ctx context.Context, domain string) (bool, error)
}

type domainPrimaryExceptionRepository struct {
	gormDb *gorm.DB
}

func NewDomainPrimaryExceptionRepository(gormDb *gorm.DB) DomainPrimaryExceptionRepository {
	return &domainPrimaryExceptionRepository{gormDb: gormDb}
}

func (d domainPrimaryExceptionRepository) Exists(ctx context.Context, domain string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainPrimaryExceptionRepository.Exists")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("domain", domain))

	var count int64
	err := d.gormDb.
		Model(&postgresentity.DomainPrimaryException{}).
		Where("domain = ?", domain).
		Count(&count).
		Error
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	return count > 0, nil
}
