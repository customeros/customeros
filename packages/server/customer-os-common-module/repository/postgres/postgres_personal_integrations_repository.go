package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/helper"
	"gorm.io/gorm"
)

type PersonalIntegrationsRepo struct {
	db *gorm.DB
}

type PersonalIntegrationRepository interface {
	FindIntegration(tenant, email, integration string) helper.QueryResult
	FindIntegrations(tenant, email string) helper.QueryResult
	SaveIntegration(integration entity.PersonalIntegration) helper.QueryResult
}

func NewPersonalIntegrationsRepo(db *gorm.DB) *PersonalIntegrationsRepo {
	return &PersonalIntegrationsRepo{db: db}
}

func (r *PersonalIntegrationsRepo) FindIntegration(tenant, email, integration string) helper.QueryResult {
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
func (r *PersonalIntegrationsRepo) FindIntegrations(tenant, email string) helper.QueryResult {
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

func (r *PersonalIntegrationsRepo) SaveIntegration(integration entity.PersonalIntegration) helper.QueryResult {
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
