package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type enrichDetailsBrandfetchRepository struct {
	db *gorm.DB
}

type EnrichDetailsBrandfetchRepository interface {
	Create(ctx context.Context, data postgres_entity.EnrichDetailsBrandfetch) (*postgres_entity.EnrichDetailsBrandfetch, error)
	GetAllSuccessByDomain(ctx context.Context, domain string) ([]postgres_entity.EnrichDetailsBrandfetch, error)
	GetLatestByDomain(ctx context.Context, domain string) (*postgres_entity.EnrichDetailsBrandfetch, error)
	GetToSyncIntoGlobalOrganizations(ctx context.Context, limit int) ([]*postgres_entity.EnrichDetailsBrandfetch, error)
	MarkSyncedToGlobalOrganizations(ctx context.Context, id uint64) error
}

func NewEnrichDetailsBrandfetchRepository(gormDb *gorm.DB) EnrichDetailsBrandfetchRepository {
	return &enrichDetailsBrandfetchRepository{db: gormDb}
}

func (r enrichDetailsBrandfetchRepository) GetAllSuccessByDomain(ctx context.Context, domain string) ([]postgres_entity.EnrichDetailsBrandfetch, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "enrichDetailsBrandfetchRepository.GetAllSuccessByDomain")
	defer spans.Finish()
	spans.LogKV("domain", domain)

	var data []postgres_entity.EnrichDetailsBrandfetch
	err := r.db.Where("domain = ? AND success = ?", domain, true).Order("created_at desc").Find(&data).Error
	if err != nil {
		return nil, err
	}

	spans.LogKV("result.count", len(data))

	return data, nil
}

func (r enrichDetailsBrandfetchRepository) GetLatestByDomain(ctx context.Context, domain string) (*postgres_entity.EnrichDetailsBrandfetch, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "enrichDetailsBrandfetchRepository.GetLatestByDomain")
	defer spans.Finish()
	spans.LogKV("domain", domain)

	var data postgres_entity.EnrichDetailsBrandfetch
	err := r.db.Where("domain = ?", domain).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spans.LogKV("result.found", false)
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	spans.LogKV("result.found", true)
	return &data, nil
}

func (r enrichDetailsBrandfetchRepository) Create(ctx context.Context, data postgres_entity.EnrichDetailsBrandfetch) (*postgres_entity.EnrichDetailsBrandfetch, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "enrichDetailsBrandfetchRepository.Create")
	defer spans.Finish()
	spans.LogObjectAsJson("data", data)

	data.CreatedAt = utils.Now()
	data.UpdatedAt = utils.Now()
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r enrichDetailsBrandfetchRepository) GetToSyncIntoGlobalOrganizations(ctx context.Context, limit int) ([]*postgres_entity.EnrichDetailsBrandfetch, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EnrichDetailsBrandfetchRepository.GetToSyncIntoGlobalOrganizations")
	defer spans.Finish()

	var data []*postgres_entity.EnrichDetailsBrandfetch
	err := r.db.
		Where("(synced_to_global_orgs IS NULL OR synced_to_global_orgs = ?) AND success = ?", false, true).
		Order("created_at asc").
		Limit(limit).
		Find(&data).Error
	if err != nil {
		return nil, err
	}

	spans.LogKV("result.count", len(data))

	return data, nil
}

func (r enrichDetailsBrandfetchRepository) MarkSyncedToGlobalOrganizations(ctx context.Context, id uint64) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "EnrichDetailsBrandfetchRepository.MarkSyncedToGlobalOrganizations")
	defer spans.Finish()

	err := r.db.Model(&postgres_entity.EnrichDetailsBrandfetch{}).Where("id = ?", id).Update("synced_to_global_orgs", true).Error
	if err != nil {
		return err
	}
	return nil
}
