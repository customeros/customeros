package postgres_repository

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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
	spans, ctx := telemetry.StartPostgresSpan(ctx, "DomainPrimaryExceptionRepository.Exists")
	defer spans.Finish()

	spans.LogKV("domain", domain)

	var count int64
	err := d.gormDb.
		Model(&postgresentity.DomainPrimaryException{}).
		Where("domain = ?", domain).
		Count(&count).
		Error
	if err != nil {
		spans.TraceError(err)
		return false, err
	}

	return count > 0, nil
}
