package postgres_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type GlobalOrganizationWebsiteToProcessRepository interface {
	GetWebsitesToProcess(ctx context.Context, limit int) ([]*postgres_entity.GlobalOrganizationWebsiteToProcess, error)
	MarkAsProcessed(ctx context.Context, id uint64, notes string) error
	AddWebsiteToProcess(ctx context.Context, website string) error
}

type globalOrganizationWebsiteToProcessRepository struct {
	db *gorm.DB
}

func NewGlobalOrganizationWebsiteToProcessRepository(gormDb *gorm.DB) GlobalOrganizationWebsiteToProcessRepository {
	return &globalOrganizationWebsiteToProcessRepository{db: gormDb}
}

func (r globalOrganizationWebsiteToProcessRepository) GetWebsitesToProcess(ctx context.Context, limit int) ([]*postgres_entity.GlobalOrganizationWebsiteToProcess, error) {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationWebsiteToProcessRepository.GetWebsitesToProcess")
	defer spans.Finish()
	spans.LogKV("limit", limit)

	var data []*postgres_entity.GlobalOrganizationWebsiteToProcess
	err := r.db.
		Where("processed IS NULL OR processed = ?", false).
		Order("created_at asc").
		Limit(limit).
		Find(&data).Error
	if err != nil {
		return nil, err
	}

	spans.LogKV("result.count", len(data))

	return data, nil
}

func (r globalOrganizationWebsiteToProcessRepository) MarkAsProcessed(ctx context.Context, id uint64, notes string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationWebsiteToProcessRepository.MarkAsProcessed")
	defer spans.Finish()
	spans.LogKV("id", id, "notes", notes)

	err := r.db.Model(&postgres_entity.GlobalOrganizationWebsiteToProcess{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed":    true,
			"processed_at": utils.Now(),
			"notes":        notes,
		}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r globalOrganizationWebsiteToProcessRepository) AddWebsiteToProcess(ctx context.Context, website string) error {
	spans, _ := telemetry.StartPostgresSpan(ctx, "GlobalOrganizationWebsiteToProcessRepository.AddWebsiteToProcess")
	defer spans.Finish()
	spans.LogKV("website", website)

	err := r.db.Create(&postgres_entity.GlobalOrganizationWebsiteToProcess{
		Website: website,
	}).Error
	if err != nil {
		return err
	}

	return nil
}
