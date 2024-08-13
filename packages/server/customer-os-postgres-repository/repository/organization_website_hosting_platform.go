package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type OrganizationWebsiteHostingPlatformRepository interface {
	GetAllUrlPatterns(ctx context.Context) ([]string, error)
}

type organizationWebsiteHostingPlatformRepository struct {
	gormDb *gorm.DB
}

func NewOrganizationWebsiteHostingPlatformRepository(gormDb *gorm.DB) OrganizationWebsiteHostingPlatformRepository {
	return &organizationWebsiteHostingPlatformRepository{gormDb: gormDb}
}

func (o organizationWebsiteHostingPlatformRepository) GetAllUrlPatterns(ctx context.Context) ([]string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationWebsiteHostingPlatformRepository.GetAllUrlPatterns")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result []string
	err := o.gormDb.Model(&entity.OrganizationWebsiteHostingPlatform{}).Pluck("url_pattern", &result).Error
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}
