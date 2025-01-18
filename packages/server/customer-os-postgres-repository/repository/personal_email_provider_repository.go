package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type PersonalEmailProviderRepository interface {
	GetPersonalEmailProviders() ([]postgres_entity.PersonalEmailProvider, error)
}

type personalEmailProviderRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewPersonalEmailProviderRepository(gormDb *gorm.DB) PersonalEmailProviderRepository {
	return &personalEmailProviderRepositoryImpl{gormDb: gormDb}
}

func (repo *personalEmailProviderRepositoryImpl) GetPersonalEmailProviders() ([]postgres_entity.PersonalEmailProvider, error) {
	span, _ := opentracing.StartSpanFromContext(context.Background(), "PersonalEmailProviderRepository.GetPersonalEmailProviders")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var result []postgres_entity.PersonalEmailProvider
	err := repo.gormDb.Find(&result).Error

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result, nil
}
