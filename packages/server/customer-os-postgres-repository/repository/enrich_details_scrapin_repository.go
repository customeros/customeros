package postgres_repository

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type enrichDetailsScrapInRepository struct {
	db *gorm.DB
}

type EnrichDetailsScrapInRepository interface {
	Create(ctx context.Context, data postgres_entity.EnrichDetailsScrapIn) (*postgres_entity.EnrichDetailsScrapIn, error)
	GetAllByParam1AndFlow(ctx context.Context, param string, flow postgres_entity.ScrapInFlow) ([]postgres_entity.EnrichDetailsScrapIn, error)
	GetLatestByParam1AndFlow(ctx context.Context, param string, flow postgres_entity.ScrapInFlow) (*postgres_entity.EnrichDetailsScrapIn, error)
	GetLatestByParam1AndFlowWithPersonFound(ctx context.Context, param string, flow postgres_entity.ScrapInFlow) (*postgres_entity.EnrichDetailsScrapIn, error)
	GetLatestByParam1AndFlowWithCompanyFound(ctx context.Context, param string, flow postgres_entity.ScrapInFlow) (*postgres_entity.EnrichDetailsScrapIn, error)
	GetLatestByAllParamsAndFlow(ctx context.Context, param1, param2, param3, param4, param5 string, flow postgres_entity.ScrapInFlow) (*postgres_entity.EnrichDetailsScrapIn, error)
	GetLatestByAllParamsAndFlowWithPersonFound(ctx context.Context, param1, param2, param3, param4, param5 string, flow postgres_entity.ScrapInFlow) (*postgres_entity.EnrichDetailsScrapIn, error)
	GetById(ctx context.Context, id uint64) (*postgres_entity.EnrichDetailsScrapIn, error)
	GetToSyncIntoGlobalOrganizations(ctx context.Context, limit int) ([]*postgres_entity.EnrichDetailsScrapIn, error)
	MarkSyncedToGlobalOrganizations(ctx context.Context, id uint64) error
}

func NewEnrichDetailsScrapInRepository(gormDb *gorm.DB) EnrichDetailsScrapInRepository {
	return &enrichDetailsScrapInRepository{db: gormDb}
}

func (r enrichDetailsScrapInRepository) GetAllByParam1AndFlow(ctx context.Context, param string, flow postgres_entity.ScrapInFlow) ([]postgres_entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetAllByParam1AndFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data []postgres_entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND flow = ?", param, flow).Find(&data).Error
	if err != nil {
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(data)))

	return data, nil
}

func (r enrichDetailsScrapInRepository) GetLatestByParam1AndFlow(ctx context.Context, param string, flow postgres_entity.ScrapInFlow) (*postgres_entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetLatestByParam1AndFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data postgres_entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND flow = ?", param, flow).Order("created_at desc").First(&data).Error
	if err != nil {
		span.LogFields(tracingLog.Bool("result.found", false))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		tracing.TraceErr(span, err)
		return nil, err // Return other errors as usual
	}

	span.LogFields(tracingLog.Bool("result.found", true))
	return &data, nil
}

func (r enrichDetailsScrapInRepository) Create(ctx context.Context, data postgres_entity.EnrichDetailsScrapIn) (*postgres_entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.Create")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	tracing.LogObjectAsJson(span, "data", data)

	data.CreatedAt = utils.Now()
	data.UpdatedAt = utils.Now()
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetById(ctx context.Context, id uint64) (*postgres_entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data postgres_entity.EnrichDetailsScrapIn
	err := r.db.Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetLatestByAllParamsAndFlow(ctx context.Context, param1, param2, param3, param4, param5 string, flow postgres_entity.ScrapInFlow) (*postgres_entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetLatestByAllParamsAndFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data postgres_entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND param2 = ? AND param3 = ? AND param4 = ? AND param5 = ? AND flow = ?", param1, param2, param3, param4, param5, flow).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetLatestByParam1AndFlowWithPersonFound(ctx context.Context, param string, flow postgres_entity.ScrapInFlow) (*postgres_entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetLatestByParam1AndFlowWithPersonFound")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data postgres_entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND flow = ? AND person_found = ?", param, flow, true).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetLatestByParam1AndFlowWithCompanyFound(ctx context.Context, param string, flow postgres_entity.ScrapInFlow) (*postgres_entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetLatestByParam1AndFlowWithCompanyFound")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data postgres_entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND flow = ? AND person_found = ?", param, flow, true).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetLatestByAllParamsAndFlowWithPersonFound(ctx context.Context, param1, param2, param3, param4, param5 string, flow postgres_entity.ScrapInFlow) (*postgres_entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetLatestByAllParamsAndFlowWithPersonFound")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	var data postgres_entity.EnrichDetailsScrapIn
	err := r.db.Where("param1 = ? AND param2 = ? AND param3 = ? AND param4 = ? AND param5 = ? AND flow = ? AND person_found = ?", param1, param2, param3, param4, param5, flow, true).Order("created_at desc").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return other errors as usual
	}

	return &data, nil
}

func (r enrichDetailsScrapInRepository) GetToSyncIntoGlobalOrganizations(ctx context.Context, limit int) ([]*postgres_entity.EnrichDetailsScrapIn, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.GetToSyncIntoGlobalOrganizations")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	// get records that has SyncedToGlobalOrgs = false or NULL, flow = ScrapInFlowCompanySearch or ScrapInFlowCompanyProfile, and CompanyFound = true
	// ordered by created date ascending, limit by limit
	var data []*postgres_entity.EnrichDetailsScrapIn
	err := r.db.
		Where("(synced_to_global_orgs IS NULL OR synced_to_global_orgs = ?) AND (flow = ? OR flow = ?) AND company_found = ?", false, postgres_entity.ScrapInFlowCompanySearch, postgres_entity.ScrapInFlowCompanyProfile, true).
		Order("created_at asc").
		Limit(limit).
		Find(&data).Error
	if err != nil {
		return nil, err
	}

	span.LogFields(tracingLog.Int("result.count", len(data)))

	return data, nil
}

func (r enrichDetailsScrapInRepository) MarkSyncedToGlobalOrganizations(ctx context.Context, id uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EnrichDetailsScrapInRepository.MarkSyncedToGlobalOrganizations")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	err := r.db.Model(&postgres_entity.EnrichDetailsScrapIn{}).Where("id = ?", id).Update("synced_to_global_orgs", true).Error
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
