package postgres_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
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
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationWebsiteToProcessRepository.GetWebsitesToProcess")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogFields(tracingLog.Int("limit", limit))

	var data []*postgres_entity.GlobalOrganizationWebsiteToProcess
	err := r.db.
		Where("processed IS NULL OR processed = ?", false).
		Order("created_at asc").
		Limit(limit).
		Find(&data).Error
	if err != nil {
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(data)))

	return data, nil
}

func (r globalOrganizationWebsiteToProcessRepository) MarkAsProcessed(ctx context.Context, id uint64, notes string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationWebsiteToProcessRepository.MarkAsProcessed")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

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
	span, _ := opentracing.StartSpanFromContext(ctx, "GlobalOrganizationWebsiteToProcessRepository.AddWebsiteToProcess")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.db.Create(&postgres_entity.GlobalOrganizationWebsiteToProcess{
		Website: website,
	}).Error
	if err != nil {
		return err
	}

	return nil
}
