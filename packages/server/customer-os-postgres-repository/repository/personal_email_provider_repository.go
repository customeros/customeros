package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type PersonalEmailProviderRepository interface {
	GetPersonalEmailProviders(ctx context.Context) ([]postgres_entity.PersonalEmailProvider, error)
}

type personalEmailProviderRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewPersonalEmailProviderRepository(gormDb *gorm.DB) PersonalEmailProviderRepository {
	return &personalEmailProviderRepositoryImpl{gormDb: gormDb}
}

func (repo *personalEmailProviderRepositoryImpl) GetPersonalEmailProviders(ctx context.Context) ([]postgres_entity.PersonalEmailProvider, error) {
	spans, ctx := telemetry.StartPostgresSpan(ctx, "PersonalEmailProviderRepository.GetPersonalEmailProviders")
	defer spans.Finish()

	var result []postgres_entity.PersonalEmailProvider
	err := repo.gormDb.Find(&result).Error

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return result, nil
}
