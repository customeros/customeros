package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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
	spans, _ := telemetry.StartPostgresSpan(ctx, "OrganizationWebsiteHostingPlatformRepository.GetAllUrlPatterns")
	defer spans.Finish()

	var result []string
	err := o.gormDb.Model(&postgres_entity.OrganizationWebsiteHostingPlatform{}).Pluck("url_pattern", &result).Error
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return result, nil
}
