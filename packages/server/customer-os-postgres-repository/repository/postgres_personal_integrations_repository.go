package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository/helper"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type PersonalIntegrationsRepo struct {
	db *gorm.DB
}

type PersonalIntegrationRepository interface {
	FindActivesByIntegration(ctx context.Context, integration string) ([]*entity.PersonalIntegration, error)
	FindIntegration(ctx context.Context, tenant, email, integration string) helper.QueryResult
	FindIntegrations(ctx context.Context, tenant, email string) helper.QueryResult
	SaveIntegration(ctx context.Context, integration entity.PersonalIntegration) helper.QueryResult
}

func NewPersonalIntegrationsRepo(db *gorm.DB) *PersonalIntegrationsRepo {
	return &PersonalIntegrationsRepo{db: db}
}

func (r *PersonalIntegrationsRepo) FindActivesByIntegration(ctx context.Context, integration string) ([]*entity.PersonalIntegration, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "PersonalIntegrationsRepo.FindActivesByIntegration")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var personalIntegrationEntity []*entity.PersonalIntegration
	err := r.db.
		Where("name = ?", integration).Where("active = ?", true).
		Find(&personalIntegrationEntity).Error

	return personalIntegrationEntity, err
}

func (r *PersonalIntegrationsRepo) FindIntegration(ctx context.Context, tenant, email, integration string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "PersonalIntegrationsRepo.FindIntegration")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var personalIntegrationEntity entity.PersonalIntegration
	err := r.db.
		Where("tenant_name = ?", tenant).
		Where("name = ?", integration).
		Where("email = ?", email).
		First(&personalIntegrationEntity).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &personalIntegrationEntity}
}

func (r *PersonalIntegrationsRepo) FindIntegrations(ctx context.Context, tenant, email string) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "PersonalIntegrationsRepo.FindIntegrations")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var personalIntegrationEntities []entity.PersonalIntegration
	err := r.db.
		Where("tenant_name = ?", tenant).
		Where("email = ?", email).
		Find(&personalIntegrationEntities).Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &personalIntegrationEntities}
}

func (r *PersonalIntegrationsRepo) SaveIntegration(ctx context.Context, integration entity.PersonalIntegration) helper.QueryResult {
	span, _ := opentracing.StartSpanFromContext(ctx, "PersonalIntegrationsRepo.SaveIntegration")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	personalIntegrationEntity := entity.PersonalIntegration{
		TenantName: integration.TenantName,
		Name:       integration.Name,
		Email:      integration.Email,
		Secret:     integration.Secret,
	}

	err := r.db.Create(&personalIntegrationEntity).Error
	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &personalIntegrationEntity}
}
